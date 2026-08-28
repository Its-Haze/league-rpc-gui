# 07: Start-with-Windows registry reconciler

**What to build:** A new `internal/startup` package that writes or removes an `HKCU\Software\Microsoft\Windows\CurrentVersion\Run` value so it matches `Behavior.LaunchAtStartup`. Reconcile on every launch. A run started by that entry opens hidden to the tray; a manual run shows the window. Use a marker the `Run` value carries (e.g. a `--hidden` arg) to tell the two apart. Task Scheduler is not used.

**Blocked by:** 01, 02

**Status:** done

- [x] `internal/startup` can set, clear, and read the `Run` value, behind an interface so tests use a fake registry
- [x] On launch the app reconciles the `Run` value to `Behavior.LaunchAtStartup` (adds if enabled, removes if not)
- [x] The `Run` command includes the marker arg; the GUI entry point starts hidden when that arg is present
- [x] A manual launch (no marker) shows the window
- [x] Toggling the setting in-app updates the registry immediately, not just on next launch
- [x] Tests against the fake registry: enable writes the value, disable removes it, reconcile is idempotent, marker arg drives hidden start

## Comments

Implemented in `internal/startup`: `RunKey` interface (real `HKCU\...\Run` backing on Windows via `golang.org/x/sys/windows/registry`, a stub off Windows, a fake in tests) and a `Reconciler` that writes only on an actual difference. `cmd/league-rpc-gui/main.go` reconciles once at launch, runs `reconcileStartupOnChange` against a `config.Store` subscription for live toggles, and creates the window with `Hidden: startup.StartedHidden(os.Args[1:])`. The `Run` command is `"<exe>" --hidden`.
