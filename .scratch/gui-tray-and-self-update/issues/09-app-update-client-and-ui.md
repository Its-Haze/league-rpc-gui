# 09: In-app update client and UI

**What to build:** Wire the Wails v3 updater service with the GitHub Releases provider and the project's ed25519 public key compiled in. Check for a newer stable release on launch and every ~6h. On a newer version, show a dismissible banner in the GUI. Clicking Update downloads the binary, verifies its signature against the `SHA256SUMS` sidecar, swaps it via the helper, and prompts to restart. Downloads start only on that click. Add a manual "Check for updates" action and a changelog view that fetches the latest release body from the GitHub API and renders it as markdown ("changelog unavailable" when offline). Inject the version at build with `-ldflags -X`; expose `App.GetVersion()`.

**Blocked by:** 01

**Status:** ready-for-agent

- [ ] Updater service configured: GitHub Releases provider, embedded ed25519 public key, stable channel only
- [ ] Launch check plus a ~6h periodic check; results surface through a binding/event
- [ ] Dismissible in-GUI banner appears only when a newer version exists
- [ ] Update flow: download on click, verify signature against `SHA256SUMS`, swap, prompt to restart
- [ ] No download traffic before the user clicks Update
- [ ] Manual "Check for updates" binding
- [ ] Changelog: fetch latest release body via the GitHub API, render markdown, degrade to "changelog unavailable" offline
- [ ] `App.GetVersion()` returns the `-ldflags`-injected version; a dev build reports a clear placeholder
- [ ] Tests: version-compare logic, "no update" and "update available" states, offline changelog fallback
