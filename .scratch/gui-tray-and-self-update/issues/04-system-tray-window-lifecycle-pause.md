# 04: System tray, window lifecycle, and the Pause flag

**What to build:** A tray icon owned by the daemon-hosting process. Left-click shows/focuses the window. Right-click menu: `Open`, `Pause presence`, `Quit`. Only `Quit` cancels the daemon context. The window's close button is intercepted and hides the window instead of quitting; the first time this happens, a one-off tray notification explains the app is still running. A single-instance guard makes a second launch surface the running window and exit. Introduce the runtime `Pause` flag (not in `Config`): default unpaused every start, clears presence immediately when set, exposed as `App.SetPaused`/`App.IsPaused`.

**Blocked by:** 01

**Status:** ready-for-agent

- [ ] Tray icon with left-click show/focus and a right-click menu of `Open` / `Pause presence` / `Quit`
- [ ] `Quit` is the only path that cancels the daemon context; `Open` and `Pause` never stop it
- [ ] Window close is redirected to hide; the app keeps running
- [ ] First hide fires a single tray notification ("still running, right-click to quit"); it does not repeat
- [ ] Single-instance guard: second launch signals the first to show its window, then exits
- [ ] `Pause` flag lives on the daemon, defaults unpaused on start, and does not persist to `Config`
- [ ] Setting `Pause` clears Discord presence immediately via the same path as League-not-running; clearing it resumes normal presence
- [ ] Tray `Pause presence` item reflects and toggles the flag
- [ ] Tests: close-hides-not-quits, pause clears presence and reset-on-restart, single-instance signalling
