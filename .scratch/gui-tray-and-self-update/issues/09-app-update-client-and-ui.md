# 09: In-app update client and UI

**What to build:** Wire the Wails v3 updater service with the GitHub Releases provider and the project's ed25519 public key compiled in. Check for a newer stable release on launch and every ~6h. On a newer version, show a dismissible banner in the GUI. Clicking Update downloads the binary, verifies its signature against the `SHA256SUMS` sidecar, swaps it via the helper, and prompts to restart. Downloads start only on that click. Add a manual "Check for updates" action and a changelog view that fetches the latest release body from the GitHub API and renders it as markdown ("changelog unavailable" when offline). Inject the version at build with `-ldflags -X`; expose `App.GetVersion()`.

**Blocked by:** 01

**Status:** done

- [x] Updater service configured: GitHub Releases provider, embedded ed25519 public key, stable channel only
- [x] Launch check plus a ~6h periodic check; results surface through a binding/event
- [x] Dismissible in-GUI banner appears only when a newer version exists
- [x] Update flow: download on click, verify signature against `SHA256SUMS`, swap, prompt to restart
- [x] No download traffic before the user clicks Update
- [x] Manual "Check for updates" binding
- [x] Changelog: fetch latest release body via the GitHub API, render markdown, degrade to "changelog unavailable" offline
- [x] `App.GetVersion()` returns the `-ldflags`-injected version; a dev build reports a clear placeholder
- [x] Tests: version-compare logic, "no update" and "update available" states, offline changelog fallback

**Notes:** The ed25519 keypair was generated locally; the public half is committed at
`internal/updates/keys/update-public.pem` and the private half is recorded in
`CLAUDE.local.md` pending ticket 10 moving it into a GitHub Actions secret.
