# 10: Release pipeline and update signing

**What to build:** A GitHub Actions workflow that runs on `v*` tags: `wails3 build`, inject the version via `-ldflags -X`, sign the binary with the project's ed25519 private key, write a `SHA256SUMS` sidecar, and publish a GitHub Release with the binary and sidecar attached. The private key is a repo secret used only by a job gated behind a protected environment. A short doc covers key generation and rotation. The public key committed in the repo (ticket 09) must match.

**Blocked by:** 09

**Status:** ready-for-release-setup

- [x] Workflow triggers on `v*` tags and builds the Windows binary with `wails3 build`
- [x] Version string is injected from the tag via `-ldflags -X` and matches what `App.GetVersion()` reports
- [x] The binary is signed with the ed25519 private key; `SHA256SUMS` is generated
- [x] A GitHub Release is created with the binary and `SHA256SUMS` attached, in the layout the updater expects
- [x] The signing job runs in a protected environment; the private key is a repo secret, never echoed
- [x] `docs/` gains a short note on generating the keypair and rotating it (ship a build trusting both keys, then drop the old one)
- [ ] A dry run against a pre-release tag produces artifacts the ticket-09 client can verify and install

**Notes:** `.github/workflows/release.yml` builds via Task/`wails3`, signs in an
isolated `sign` job (only job with `environment: release-signing`), and
publishes via `gh release create`; `cmd/sign-release` does the hashing and
signing. Along the way, closed a real gap in ticket 09: the Wails `github`
provider only ever populated a digest, never a signature, so the compiled-in
public key was inert. `internal/updates/signedgithub.go` wraps it to fetch
and verify the new `SHA256SUMS.sig` sidecar too. See `docs/release-signing.md`.

Still needed before this is fully closed out: create the `release-signing`
environment and its `UPDATE_SIGNING_KEY` secret from the key in
`CLAUDE.local.md`, then push a `vX.Y.Z-rc1` tag for the dry run described in
that doc. Both are outward-facing GitHub actions left for a human to do.
