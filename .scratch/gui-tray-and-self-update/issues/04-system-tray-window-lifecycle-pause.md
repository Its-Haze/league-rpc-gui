# 04: System tray, window lifecycle, and the Pause flag

**What to build:** A tray icon owned by the daemon-hosting process. Left-click shows/focuses the window. Right-click menu: `Open`, `Pause presence`, `Quit`. Only `Quit` cancels the daemon context. The window's close button is intercepted and hides the window instead of quitting; the first time this happens, a one-off tray notification explains the app is still running. A single-instance guard makes a second launch surface the running window and exit. Introduce the runtime `Pause` flag (not in `Config`): default unpaused every start, clears presence immediately when set, exposed as `App.SetPaused`/`App.IsPaused`.

**Blocked by:** 01

**Status:** done

- [x] Tray icon with left-click show/focus and a right-click menu of `Open` / `Pause presence` / `Quit` (`cmd/league-rpc-gui/main.go`, `tray.go`)
- [x] `Quit` is the only path that cancels the daemon context; `Open` and `Pause` never stop it (`Quit` calls `wailsApp.Quit()` -> `OnShutdown` -> `cancel()`)
- [x] The daemon goroutine starts only after `application.New` has run the single-instance guard, so a second launch exits before it can touch Discord
- [x] Window close is redirected to hide; the app keeps running (`RegisterHook(events.Common.WindowClosing)` -> `handleClose()` + `e.Cancel()`)
- [x] First hide shows a single one-off message; it does not repeat (`trayController.firstHide sync.Once`). Uses a Wails info dialog, not an OS toast: the notifications service needs a per-launch `HKCU\...\CLSID` registry write and a new `go-toast` dependency, too much for a "still running" hint
- [x] Single-instance guard: second launch signals the first to show its window, then exits (`application.SingleInstanceOptions`, `OnSecondInstanceLaunch` shows/focuses the window)
- [x] `Pause` flag lives on the daemon, defaults unpaused on start, and does not persist to `Config` (`Daemon.paused atomic.Bool`, zero value = unpaused, never read/written by `config`)
- [x] Setting `Pause` clears Discord presence immediately via the same path as League-not-running; clearing it resumes normal presence (`pauseSignal` wakes `presenceLoop.reconcile`, which clears on the same branch shape as no-League)
- [x] Tray `Pause presence` item reflects and toggles the flag. Both the tray click and the frontend `guiService.SetPaused` binding route through `trayController.setPaused`, which updates the daemon flag and mirrors it onto the checkbox (`reflectChecked`, marshalled to the main thread), so the two never disagree. Also exposed as `App.SetPaused`/`App.IsPaused`
- [x] Tests: close-hides-not-quits + pause-toggle (`cmd/league-rpc-gui/tray_test.go`), pause clears presence / resumes / defaults unpaused (`internal/daemon/daemon_test.go`), pause delegates through `App` (`internal/app/app_test.go`). Single-instance signalling and the real close-to-tray are Wails runtime behavior, verified by `wails3 build` green and left for a manual desktop check
