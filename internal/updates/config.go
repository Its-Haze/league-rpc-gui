package updates

import (
	_ "embed"
	"time"

	wupdater "github.com/wailsapp/wails/v3/pkg/updater"
	"github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

// RepoSlug is the GitHub repository the App Update flow pulls releases from.
const RepoSlug = "Its-Haze/league-rpc-gui"

// ChecksumAsset is the sidecar every release publishes alongside the binary.
// Its ed25519 signature is verified against publicKeyPEM before any swap.
const ChecksumAsset = "SHA256SUMS"

// CheckInterval is how often the background check re-hits GitHub Releases
// after the launch check.
const CheckInterval = 6 * time.Hour

// publicKeyPEM is the release-signature trust root, fixed at build time. It
// is the ONLY authority for verifying a release's SHA256SUMS signature: the
// matching private key lives outside this repo (see CLAUDE.local.md until
// ticket 10 moves it into a GitHub Actions secret) and signs each tagged
// release build. Rotating the key means shipping a build that trusts both
// the old and new key before the old one is retired.
//
//go:embed keys/update-public.pem
var publicKeyPEM []byte

// BuildConfig assembles the Wails updater configuration for currentVersion:
// GitHub Releases as the only provider, stable channel, no built-in window.
func BuildConfig(currentVersion string) (wupdater.Config, error) {
	provider, err := github.New(github.Config{
		Repository:    RepoSlug,
		Prerelease:    false,
		ChecksumAsset: ChecksumAsset,
	})
	if err != nil {
		return wupdater.Config{}, err
	}
	return wupdater.Config{
		CurrentVersion: currentVersion,
		Providers:      []wupdater.Provider{provider},
		PublicKey:      publicKeyPEM,
		Window:         wupdater.WindowNone,
	}, nil
}
