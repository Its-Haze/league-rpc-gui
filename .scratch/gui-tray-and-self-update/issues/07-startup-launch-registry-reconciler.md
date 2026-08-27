# 07: Start-with-Windows registry reconciler

**What to build:** A new `internal/startup` package that writes or removes an `HKCU\Software\Microsoft\Windows\CurrentVersion\Run` value so it matches `Behavior.LaunchAtStartup`. Reconcile on every launch. A run started by that entry opens hidden to the tray; a manual run shows the window. Use a marker the `Run` value carries (e.g. a `--hidden` arg) to tell the two apart. Task Scheduler is not used.

**Blocked by:** 01, 02

**Status:** ready-for-agent

- [ ] `internal/startup` can set, clear, and read the `Run` value, behind an interface so tests use a fake registry
- [ ] On launch the app reconciles the `Run` value to `Behavior.LaunchAtStartup` (adds if enabled, removes if not)
- [ ] The `Run` command includes the marker arg; the GUI entry point starts hidden when that arg is present
- [ ] A manual launch (no marker) shows the window
- [ ] Toggling the setting in-app updates the registry immediately, not just on next launch
- [ ] Tests against the fake registry: enable writes the value, disable removes it, reconcile is idempotent, marker arg drives hidden start
