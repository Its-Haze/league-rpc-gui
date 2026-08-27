# Daemon skeleton: always-on background presence

Status: ready-for-agent

## Problem Statement

Today the rewrite has no runnable entry point at all. `cmd/league-rpc/` is empty. The eventual target (confirmed during grilling) is a Blitz-style background app: launched at Windows startup or manually, living in the system tray, that a player never has to manually restart. It should keep working across every ordinary event in a play session, Discord not open yet, League not open yet, League closing and reopening, Discord closing and reopening, without the player ever touching it again. Right now none of that exists: there's no process that stays alive, no logic that tolerates League or Discord being unavailable, and no logic that recovers when they come back.

## Solution

Build the Daemon skeleton: a long-running process that starts two independent Connection Supervisors (one for Discord, one for the LCU), never exits on its own (only on crash or an explicit shutdown signal), and drives the existing presence-mapping logic (`discord.MapStateToPresence`, `discord.ShouldClearPresence`) for every gameflow phase except the in-game detail phase. This proves the daemon's core lifecycle contract (stays alive, reconnects, shows correct presence, clears correctly) before champion/skin detail, spectator mode, and the other feature-parity work (tracked separately) gets layered on top of it.

## User Stories

1. As a player, I want the Daemon to launch and immediately start trying to connect to Discord and to League, so that I don't have to manually trigger anything for presence to start working.
2. As a player, I want the Daemon to keep running even when League isn't open yet, so that I can leave it running in the background before I've launched League.
3. As a player, I want the Daemon to keep running even when Discord isn't open yet, so that I don't have to worry about launch order between Discord and League.
4. As a player, I want my Discord presence to appear automatically once both Discord and League are connected, so that I don't have to trigger anything manually once both apps are up.
5. As a player, I want my Discord presence to disappear immediately when I close League, so that I'm not showing a stale or incorrect status to friends.
6. As a player, I want the Daemon to reconnect to League automatically if I close and reopen it, so that I don't have to restart the Daemon every time I relaunch League.
7. As a player, I want the Daemon to reconnect to Discord automatically if I close and reopen Discord, so that presence resumes without restarting the Daemon.
8. As a player, I want the Daemon to never crash or exit just because League or Discord isn't running, so that it behaves like a real background service rather than a script that needs supervision.
9. As a player, I want my Discord presence to reflect that I'm idle in the League client, so that friends can see I'm in the client even before I queue up.
10. As a player, I want my Discord presence to reflect that I'm in a lobby, so that friends know I'm setting up a game.
11. As a player, I want my Discord presence to reflect that I'm in queue/matchmaking, so that friends know I'm searching for a match.
12. As a player, I want my Discord presence to reflect that I'm in champion select, so that friends know I'm picking my champion.
13. As a player, I want my Discord presence to keep displaying correctly even if League's own native Discord integration tries to take over, so that my chosen presence isn't silently replaced.
14. As a player, I want to see a "Launching League..." placeholder presence once I've actually started League, so that I'm not showing nothing while I wait for the client to finish loading.
15. As a player, I want that placeholder to disappear the moment League's client is actually ready, so that it doesn't linger and show stale "launching" text once I'm in the client.
16. As a player, I want no placeholder and no presence at all while League simply isn't running, so that the Daemon never shows a misleading "launching" card when I haven't even started the game.
17. As a player, I want the Daemon to keep running unattended for arbitrarily long periods (days, through sleep/wake cycles, through many separate League sessions) without needing a restart, so that it behaves like the "set it and forget it" background apps I'm used to.
18. As a developer, I want a single top-level entry point that starts the Discord and LCU Connection Supervisors and blocks until an explicit shutdown signal, so that the Daemon's lifecycle contract is easy to reason about and test.
19. As a developer, I want the Discord Connection Supervisor and the LCU Connection Supervisor to run independently of each other, so that one being unavailable never blocks or delays the other from connecting.
20. As a developer, I want Discord-process detection (waiting for Discord's process to exist before attempting an IPC login) implemented via `gopsutil`, so that the Daemon doesn't fail outright the moment it starts if Discord hasn't launched yet.
21. As a developer, I want League-*process* detection (also via `gopsutil`, same mechanism as Discord's) checked independently of the LCU connection itself, so that the Daemon can distinguish "League process running, LCU not connected yet" (show the placeholder) from "League process not running at all" (show nothing).
22. As a developer, I want the actual LCU connection handshake to continue relying on `lcu-gopher`'s existing `AwaitConnection` behavior rather than a second, duplicate connect-polling loop, so that "wait for the LCU API to come up" logic isn't implemented twice.
23. As a developer, I want the existing debounced `Updater` extended to keep resending the current presence periodically (heartbeat) and in a fast burst right after a real change (reclaim), so that League's native presence integration can't silently overwrite ours (per ADR-0001).
24. As a developer, I want every Discord presence send, heartbeat, reclaim, placeholder rotation, and real updates alike, to go through the same mutex-guarded `Updater`, so that no two goroutines can race on the same Discord IPC connection.
25. As a developer, I want the Daemon's presence-clearing behavior to be uniform with exactly one exception (no LCU connection means no presence, with no grace period, except while the League process is detected running and the LCU hasn't connected yet, which shows the placeholder instead), so that behavior matches the refined ADR-0002 exactly.
26. As a developer, I want the Daemon skeleton fully testable without a real Discord client or a real League client running, so that CI and local test runs never depend on external processes.
27. As a developer, I want champion/skin/chroma detection, spectator mode, League Classic, ARAM Mayhem, and locale support explicitly excluded from this spec, so that this milestone stays scoped to "the Daemon runs forever and shows correct presence for every phase except in-game detail."
28. As a developer, I want this milestone to ship with no GUI and no system tray icon, so that autostart-at-Windows-login and tray behavior can be built later against a Daemon core that's already proven to run correctly, rather than being built and debugged at the same time as the core lifecycle logic.

## Implementation Decisions

- **New package `internal/daemon`**: a `Daemon` type whose dependencies (a Discord connector, an LCU connector, a clock/ticker abstraction) are injected rather than constructed internally. It exposes one method, `Run(ctx context.Context) error`, which starts both Connection Supervisors and blocks until `ctx` is canceled. This is the single top-level seam for this spec (see Testing Decisions).
- **Connection Supervisor** (a pattern applied twice, once per side, not necessarily one shared type): retries the underlying `Connect()` indefinitely on a short fixed interval (a few seconds, no exponential backoff needed, since these are local process/IPC checks rather than calls to a rate-limited remote service). A supervisor never returns/exits on failure; failure just means "try again."
  - **Discord Connection Supervisor**: checks `internal/process` for a running Discord process before each `discord.Client.Connect()` attempt. On a detected IPC disconnect (e.g. a failed send from `Updater`'s heartbeat), it resets to "wait for process" and retries from there.
  - **LCU Connection Supervisor**: wraps the existing `lcu.Client.Connect()`/`Disconnect()`. Since `lcu.Client.Connect()` already blocks awaiting League via `AwaitConnection: true`, the supervisor's loop is simply: call `Connect()` (blocks until League is up), run until a disconnect is detected, call `Disconnect()` to clean up, loop back to `Connect()`.
- **New package `internal/process`**: exposes a check for whether a named process is currently running, using `gopsutil`. Used for two different name lists: Discord (Discord, Discord Canary, etc., the same configurable list the Python original used) by the Discord Connection Supervisor, and League (`RiotClientServices.exe`/`LeagueClientUx.exe`, matching the Python reference's own detection) to drive the placeholder trigger below. Both are single checks called repeatedly by their retry loops, not long-lived blocking waits.
- **Placeholder presence trigger**: the Daemon tracks League-process-detected (via `internal/process`) as a signal independent of the LCU Connection Supervisor's own connect state. This distinguishes two states that both look like "LCU not connected": League's process running with the LCU not yet up (show the rotating "Launching League..." placeholder, ported from the Python reference), versus League's process not running at all (show nothing, per below). The placeholder stops the instant the LCU Connection Supervisor reports connected; real presence takes over from there.
- **`discord.Updater` (existing type, `internal/discord/updater.go`) is extended, not replaced**, per ADR-0001: its one-shot `time.AfterFunc` debounce becomes a repeating ticker that continues firing (the heartbeat) after the first real send, with a shorter interval for a few cycles right after a real change (the reclaim burst). The existing `sync.Mutex` on `Updater` is reused as the single point of serialization; no new lock is introduced, and `discord.Client.UpdatePresence`/`ClearPresence` remain the only entry points into the underlying `rich-go` client, called exclusively from within `Updater`.
- **Presence contract** (refined ADR-0002): no LCU connection means no presence, with no grace period, except while League's process is detected running and the LCU hasn't connected yet, which shows the placeholder instead of nothing. Both "never connected" and "just disconnected" behave identically once that one exception is accounted for.
- **`Daemon.Run` shutdown sequence**: on `ctx` cancellation, both supervisors are stopped, presence is cleared, and both underlying clients are disconnected (best-effort, matching the existing "never let cleanup panic on the way out" posture already present in the Python reference) before `Run` returns.
- Presence mapping itself (`discord.MapStateToPresence`, `discord.ShouldClearPresence`, `internal/state`) is **not modified** by this spec beyond what's needed to wire it into `Daemon`/`Updater`. It already covers every in-scope `GameFlowPhase` (`None`/in-client, `Lobby`, `Matchmaking`/`ReadyCheck`, `ChampSelect`/`GameStart`). `GameFlowInProgress` continues to route through the existing `BuildInGamePresence`/`BuildTFTInGamePresence` builders as-is; this spec does not add champion/skin detail to them (see Out of Scope).
- The four TFT queue-variant constants decided in the grilling session (Ranked 1100, Double Up 1160, Hyper Roll 1130, alongside the existing normal-TFT 1090) may be added to `pkg/types/queue.go` as part of this work since they're cheap and already decided, but nothing in this spec consumes them yet.

## Testing Decisions

- There is no existing test suite in `league-rpc/` yet (`**/*_test.go` returned nothing). There's no prior art in this codebase to follow; use standard Go table-driven test conventions.
- **Primary seam**: `daemon.Daemon.Run`. Tests inject fake Discord/LCU connectors (satisfying minimal `Connect()`/`Disconnect()` interfaces) and a fake clock, and assert:
  - `Run` does not return until `ctx` is canceled, even when both fake connectors fail repeatedly.
  - Once both fakes report connected, a presence update is issued.
  - When the fake LCU connector reports disconnected *and* the fake League-process check reports "not running," presence is cleared immediately, with no delay.
  - When the fake LCU connector reports disconnected *but* the fake League-process check reports "running," the placeholder is shown instead of nothing.
  - The placeholder stops and real presence takes over the moment the fake LCU connector reports connected.
  - After a simulated LCU disconnect, the supervisor retries `Connect()` again without `Run` returning an error.
  - The Discord and LCU supervisors' failures are independent: one failing repeatedly doesn't block the other from succeeding.
- **Secondary seam (already exists, unchanged by this spec)**: `discord.MapStateToPresence`/`ShouldClearPresence` are pure functions over `state.State` + `config.Config`, table-driven tests over one `state.State` fixture per in-scope `GameFlowPhase`, asserting the resulting `RPCData` shape. This spec doesn't change their signature, only ensures `Daemon` calls them at the right time.
- `internal/process`'s Discord- and League-detection checks should each be tested against a mocked process-listing call, not a real Discord/League process: CI has neither installed.
- `Updater`'s heartbeat/reclaim timing must be tested via the injected clock/ticker abstraction, not real `time.Sleep`/wall-clock waits. Timing tests should run in milliseconds regardless of the real heartbeat interval.
- Good tests here assert observable behavior (was a presence update issued, was presence cleared, did `Run` block or return), not internal implementation details like which goroutine did the work or exact retry counts, since those are free to change as the supervisor loop is tuned.

## Out of Scope

- Champion/skin/chroma detection and any in-game presence detail (`GameFlowInProgress` beyond what already exists).
- Spectator mode (`Watching` phase).
- League Classic, ARAM Mayhem, and any other new game-mode/queue display work.
- Locale/language threading through champion/skin lookups (locale detection itself, and everything that consumes it, is part of the champion/skin work this spec excludes).
- The `internal/livegame` (Live Client Data) and `internal/championdata` (Data Dragon/Community Dragon) packages decided during grilling: neither is needed until the in-game-detail work begins.
- Auto-launching League on the Daemon's behalf. `config.AutoLaunchLeague` already exists in `internal/config` but its behavior is untouched here. The placeholder trigger (League process detected) works identically whether the player launched League themselves or the Daemon eventually does it for them.
- Any Wails GUI, system tray icon, or Windows-startup (autostart) registration. The Daemon must be *capable* of running indefinitely as a background process; wiring it into a tray shell and an autostart entry is a separate, later packaging milestone.
- Any frontend/React work.

## Further Notes

- This spec is the first of at least two: it proves the Daemon's core lifecycle (stays alive, reconnects independently on both sides, shows correct presence outside of in-game detail, clears correctly) before the feature-parity work identified from the upstream `league-rpc-linux` diff (champion/skin fixes, spectator mode, League Classic, ARAM Mayhem, locale) gets built on top of it as a separate spec.
- Builds directly on [ADR-0001](../docs/adr/0001-updater-owns-the-heartbeat.md) (heartbeat/reclaim folded into `Updater`) and [ADR-0002](../docs/adr/0002-daemon-never-exits-on-connection-loss.md) (the Daemon never exits on connection loss; no LCU connection means no presence, except while League's process is detected running and the LCU hasn't connected yet).
- Confirmed: single seam. `Daemon.Run` is the only point tests attach to; no separate seams per Connection Supervisor. If precise backoff-timing assertions turn out to be awkward to express through the black-box `Run` seam later, revisit this: the `Connect()`/`Disconnect()` interfaces on each supervisor are already narrow enough to become their own seam without restructuring anything, but that's not needed for this milestone.
