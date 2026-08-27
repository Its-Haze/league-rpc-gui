# 10: Release pipeline and update signing

**What to build:** A GitHub Actions workflow that runs on `v*` tags: `wails3 build`, inject the version via `-ldflags -X`, sign the binary with the project's ed25519 private key, write a `SHA256SUMS` sidecar, and publish a GitHub Release with the binary and sidecar attached. The private key is a repo secret used only by a job gated behind a protected environment. A short doc covers key generation and rotation. The public key committed in the repo (ticket 09) must match.

**Blocked by:** 09

**Status:** ready-for-agent

- [ ] Workflow triggers on `v*` tags and builds the Windows binary with `wails3 build`
- [ ] Version string is injected from the tag via `-ldflags -X` and matches what `App.GetVersion()` reports
- [ ] The binary is signed with the ed25519 private key; `SHA256SUMS` is generated
- [ ] A GitHub Release is created with the binary and `SHA256SUMS` attached, in the layout the updater expects
- [ ] The signing job runs in a protected environment; the private key is a repo secret, never echoed
- [ ] `docs/` gains a short note on generating the keypair and rotating it (ship a build trusting both keys, then drop the old one)
- [ ] A dry run against a pre-release tag produces artifacts the ticket-09 client can verify and install
