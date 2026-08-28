// Command sign-release is a release-pipeline tool, not a shipped artifact.
package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const signingKeyEnv = "UPDATE_SIGNING_KEY"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "sign-release:", err)
		os.Exit(1)
	}
}

func run() error {
	artifact := flag.String("artifact", "", "path to the release binary to hash and sign")
	sumsOut := flag.String("sums-out", "SHA256SUMS", "path to write the sha256sum-style digest listing")
	sigOut := flag.String("sig-out", "SHA256SUMS.sig", "path to write the detached ed25519 signature listing")
	flag.Parse()

	if *artifact == "" {
		return fmt.Errorf("-artifact is required")
	}

	priv, err := loadPrivateKey()
	if err != nil {
		return err
	}

	digest, err := hashFile(*artifact)
	if err != nil {
		return fmt.Errorf("hash %s: %w", *artifact, err)
	}
	name := filepath.Base(*artifact)

	sig := ed25519.Sign(priv, digest)

	if err := writeLine(*sumsOut, hex.EncodeToString(digest), name); err != nil {
		return err
	}
	if err := writeLine(*sigOut, hex.EncodeToString(sig), name); err != nil {
		return err
	}

	fmt.Printf("signed %s (sha256 %s)\n", name, hex.EncodeToString(digest))
	return nil
}

// loadPrivateKey reads and parses the PKCS#8 PEM ed25519 key from the
// signing-key environment variable. It never appears in an error message.
func loadPrivateKey() (ed25519.PrivateKey, error) {
	raw := os.Getenv(signingKeyEnv)
	if raw == "" {
		return nil, fmt.Errorf("%s is not set", signingKeyEnv)
	}
	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return nil, fmt.Errorf("%s does not contain a PEM block", signingKeyEnv)
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}
	priv, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is %T, want ed25519.PrivateKey", key)
	}
	return priv, nil
}

func hashFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// writeLine writes a single sha256sum-style line: "<hex>  <filename>".
func writeLine(path, hexValue, filename string) error {
	return os.WriteFile(path, []byte(hexValue+"  "+filename+"\n"), 0o644)
}
