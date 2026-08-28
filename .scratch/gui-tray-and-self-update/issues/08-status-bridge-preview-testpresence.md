# 08: Status bridge, last-sent preview, and Test presence

**What to build:** A bridge in `internal/app` that subscribes to `state.Manager.Updates()` and reads the supervisors' `Connected()` / `LeagueProcessDetected()` plus the `Pause` flag, assembles a `StatusSnapshot` (League process, LCU, and Discord connection states; current `GameFlowPhase`; the presence the `Updater` last sent), exposes `GetStatus()`, and emits `status:changed` when it changes. Add `Updater.LastSent()` returning the most recent presence payload. Add `App.TestPresence()` that pushes a fixed sample presence through the real `Updater` for ~30s and then reverts.

**Blocked by:** 01

**Status:** done

- [x] `StatusSnapshot` type with the three connection states, `GameFlowPhase`, and last-sent presence
- [x] Bridge assembles it from `state.Manager` + supervisor accessors + `Pause`, with no direct daemon coupling beyond those interfaces
- [x] `GetStatus()` binding returns the current snapshot; `status:changed` fires on change and not on a no-op update
- [x] `Updater.LastSent()` returns the last presence it sent (or a cleared marker), safe for concurrent read
- [x] The preview surface reads `LastSent()`, never a recomputation from `State` + `Config`
- [x] `App.TestPresence()` shows the sample for ~30s then returns to live state, and is a no-op-safe if called again mid-window
- [x] Tests: snapshot shape, change detection, `LastSent` after a send and after a clear, `TestPresence` reverts

## Implementation notes

- `internal/app/status.go`: `statusBridge` + `Connections` / `PresenceProbe` / `TestPresenter` interfaces. `App` gains `Option`s (`WithStatus`), `GetStatus`, `OnStatusChange`, `RunStatus`, `TestPresence`.
- `state.Manager.Subscribe()`: independent fan-out channel so the bridge and the presence loop don't steal each other's events off the single `Updates()` channel.
- `discord.Updater.LastSent()` + `LastSent` struct; `PushSample` sends a raw payload and defends it with the heartbeat. `discord.BuildTestPresence()` is the fixed sample.
- `daemon`: `TestPresence()`, a 30s `testActive` window `presenceLoop` honors (live updates, config resends, placeholders all suppressed; pause still wins), `testSignal` re-sends live presence when the window closes. Accessors `DiscordConnected` / `LCUConnected` / `LeagueProcessDetected` / `LastSent` / `SubscribeState`.
- Frontend Wails bindings for `GetStatus` / `TestPresence` / `StatusSnapshot` regenerate with the frontend work (tickets 11-16).
