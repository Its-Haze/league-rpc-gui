# GUI, system tray, and in-app self-update

Status: ready-for-agent

## Problem Statement

The daemon runs headless. There is no way for a user to see what it is doing, change what the presence shows, make it start with Windows, or update it without downloading a new binary from GitHub by hand. The backend was built for this: `internal/app` already exists as "the seam the GUI binds to" with plain-Go methods, `config.Store` already has a `Subscribe()` channel for live reload, and `daemon.Wire()` carries a comment saying the GUI and the headless launcher both call it. Nothing consumes any of that yet.

This milestone builds the GUI (Wails v3 + React + TypeScript), the system tray behavior, the config surface that lets a user shape the presence, a live status view, and an in-app updater that pulls signed releases from GitHub.

## Solution

One binary hosts everything (see [ADR-0004](../../docs/adr/0004-gui-process-hosts-the-daemon.md)): the GUI process starts `daemon.Wire()` in a goroutine and runs the Wails event loop. Closing the window hides it to the tray; only the tray "Quit" ends the process. Windows starts the app at logon through an `HKCU\...\Run` value the app reconciles to a config setting.

`Config` grows from a flat struct into a nested tree with a schema version and a one-time migration. Per-`GameMode` overrides layer on top of global display defaults. Presence text becomes user-editable per-context templates rendered by a single Go engine that the GUI previews through a binding.

The GUI reads a `StatusSnapshot` assembled from `State` and the supervisors' accessors, pushed on change. The live preview shows the presence the `Updater` last sent, never a recomputation.

The updater ([ADR-0005](../../docs/adr/0005-in-app-update-via-signed-binary-swap.md)) checks GitHub Releases, and on user click downloads the new binary, verifies its ed25519 signature against a `SHA256SUMS` sidecar, swaps it, and restarts. A tagged CI build produces and signs those artifacts.

## User Stories

1. As a user, I want a window that shows whether League, the League Client API, and Discord are each connected right now, so that I can tell at a glance whether presence is working.
2. As a user, I want to see a preview of exactly what my Discord presence currently looks like, so that I don't have to alt-tab to Discord to check.
3. As a user, I want to toggle whether my rank, my KDA/CS, and the status emojis appear in presence, so that I control how much detail friends see.
4. As a user, I want those toggles to differ per game mode, so that I can show my rank in Ranked but hide it in Arena.
5. As a user, I want to edit the text of my presence for each phase (in client, champ select, in game, spectating), so that it reads the way I want.
6. As a user, I want a live preview of my edited template against sample data, so that I can see the result before saving.
7. As a user, I want a setting that makes the app start when Windows starts, so that I never have to launch it manually.
8. As a user, I want the app to start hidden in the tray when Windows launches it, so that it doesn't interrupt my login.
9. As a user, I want clicking the window's "X" to hide the app to the tray and keep presence running, so that closing the window doesn't stop the app.
10. As a user, I want to be told, the first time the window hides, that the app is still running, so that I don't think I quit it.
11. As a user, I want to click the tray icon to bring the window back, so that reopening is one action.
12. As a user, I want a tray right-click menu with Open, Pause presence, and Quit, so that I can control the app without the window.
13. As a user, I want "Pause presence" to clear my Discord presence immediately, so that I can go invisible without quitting.
14. As a user, I want pause to reset itself when the app restarts, so that I can't accidentally leave presence off forever.
15. As a user, I want a "Test presence" button that shows a sample presence in Discord for a few seconds, so that I can see my settings take effect without launching League.
16. As a user, I want a first-run walkthrough that explains what the app does and checks that Discord and League are reachable, so that setup is obvious.
17. As a user, I want a Help section with a log viewer, a button to open the logs folder, and a button to copy diagnostics, so that I can get help when something breaks.
18. As a user, I want links to the Discord community server and to file a bug or feature request, so that I know where to go.
19. As a user, I want the app to tell me when a new version is out, so that I don't have to check GitHub.
20. As a user, I want to update to that version from inside the app, so that I never download and replace a binary by hand.
21. As a user, I want update downloads to start only when I click Update, so that the app doesn't use my bandwidth for updates I didn't ask for.
22. As a user, I want to see the changelog for the new version before I update, so that I know what changed.
23. As a user, I want a light, dark, and system theme, so that the app matches the rest of my desktop.
24. As a developer, I want one binary that hosts both the GUI and the daemon with no IPC, so that live config and status need no second transport.
25. As a developer, I want the existing `cmd/league-rpc` headless build kept for debugging, so that I can run the daemon without the GUI.
26. As a developer, I want `Config` migrated to a versioned nested tree with a one-time upgrade from the current flat file, so that per-mode overrides can tell "set false" from "not set."
27. As a developer, I want one presence template engine in Go used for both real sends and the GUI preview, so that the two can't drift.
28. As a developer, I want logging to write to a rotating file and an in-memory ring buffer, so that a GUI process with no console still produces logs and can show them live.
29. As a developer, I want the status read model assembled in one place and pushed to the GUI on change, so that screens don't each poll the daemon.
30. As a developer, I want the live preview to be the presence the `Updater` last sent, so that it can't disagree with what Discord shows.
31. As a developer, I want release artifacts (binary, `SHA256SUMS`, ed25519 signature) produced and signed by a tagged CI build, so that the updater has something to verify against.
32. As a developer, I want the frontend's logic-bearing parts (template preview, config form state, override resolution) covered by Vitest, so that regressions in those are caught.

## Implementation Decisions

- **New entry point `cmd/league-rpc-gui/main.go`**: constructs the logger, `config.Store`, calls `daemon.Wire(store, logger)`, runs `Daemon.Run` in a goroutine, and starts the Wails v3 app with the `internal/app` bindings and a tray. The Wails event loop owns the main thread. `cmd/league-rpc` stays as-is for headless runs and is not a shipped artifact.
- **Frontend layout**: `frontend/` hand-assembled (Vite + React + TypeScript + Tailwind v4 + Radix primitives), not from `wails3 init`. No component-library dependency; components are built in-repo.
- **`internal/app` stays Wails-free**: bindings are plain Go. A thin adapter in `cmd/league-rpc-gui` bridges `internal/app` method calls and channel subscriptions to Wails events.
- **Config schema v2** (`internal/config`): nested tree, top-level `SchemaVersion int`, `DiscordAppID`, `Theme`. Groups: `Display` (`Default.{ShowRank,ShowStats}`, `Modes map[string]ModeOverride`), `Presence` (`ShowEmojis`, `ShowInClient`, `Idle`, `Templates map[string]TemplatePair`), `Behavior` (`LaunchAtStartup`), `Advanced` (`UpdateInterval`, `StatsPollingInterval`, `DebugMode`). `DiscordAppIDPresets` unchanged. `Load` detects a versionless (flat) file, migrates it into the tree, writes it back. `clamp()` and `DefaultConfig()` gain a branch/default per new field. All daemon consumers (`discord/mapper.go` and anything else reading `config.Config`) move to the new paths.
- **Mode key**: the `GameMode` display concept from `CONTEXT.md`, not a queue ID. Only `ShowRank` and `ShowStats` are per-mode. Resolution helper: a mode with no override entry inherits `Display.Default`.
- **Presence template engine** (`internal/presence/template`): plain `{token}` substitution, no logic, no Go `text/template`. Per-context token sets for `in-client`, `champ-select`, `in-game`, `spectating` (spectating = the `Watching` phase). Unknown token: left literal, reported by the preview API. Missing data at send time: token resolves to empty, then runs of whitespace and dangling separators (`-`, `|`, `•`) collapse. Ships default templates that reproduce the current hardcoded presence strings exactly. `discord/mapper.go` renders through this engine.
- **Runtime `Pause`**: a flag on the daemon (not `Config`), default unpaused on every start. While paused, presence clears immediately, same path as League-not-running. Exposed as `App.SetPaused(bool)` / `App.IsPaused()` and surfaced in the tray menu and the top strip.
- **Status bridge** (`internal/app`): subscribes to `state.Manager.Updates()`, reads `discordRunner.Connected()`, `lcuRunner.Connected()`, `lcuRunner.LeagueProcessDetected()`, and the paused flag, assembles a `StatusSnapshot` (three connection states, `GameFlowPhase`, last-sent presence). Exposes `GetStatus()` and emits `status:changed` on change. `Updater` gains a `LastSent()` accessor; the preview reads that, never a recomputation.
- **`App.TestPresence()`**: pushes a fixed sample presence through the real `Updater` for ~30s, then reverts to live state.
- **Startup launch** (`internal/startup`): writes/removes `HKCU\Software\Microsoft\Windows\CurrentVersion\Run` to match `Behavior.LaunchAtStartup`, reconciled on every launch. A run launched by that entry starts hidden to the tray; a manual run shows the window. Mechanism: an arg (e.g. `--hidden`) the `Run` value includes, or a first-run heuristic if that proves cleaner. Task Scheduler is not used.
- **Window lifecycle**: window close intercepted and redirected to hide. First hide fires a one-off tray notification. Single-instance guard: a second launch signals the running instance to show its window and exits.
- **Logging** (`internal/logging`): zerolog multi-writer to a `lumberjack`-rotated file at `%APPDATA%\league-rpc\logs\league-rpc.log` (size-capped, keep 2) and an in-memory ring buffer (~500 lines). Level follows `Advanced.DebugMode`. `GetRecentLogs()` binding returns the buffer; new lines emit as `log:line` events. `cmd/league-rpc` keeps console output.
- **Updater client**: Wails v3 updater service, GitHub Releases provider, ed25519 public key compiled in. Check on launch and every ~6h. Stable channel only. Dismissible in-GUI banner on a newer version. Download starts on user click: download, verify signature against `SHA256SUMS`, swap via the helper, prompt to restart. Manual "Check for updates" in About. About fetches the latest release body from the GitHub API and renders it as markdown; offline shows "changelog unavailable".
- **Version**: injected at build with `-ldflags -X` from the git tag; `App.GetVersion()` binding.
- **Release pipeline**: GitHub Actions on `v*` tags runs `wails3 build`, signs the binary with an ed25519 key held as a repo secret behind a protected environment, writes `SHA256SUMS`, and publishes a GitHub Release with the binary and sidecar. Semver tags.
- **GitHub issue templates**: a bug template and a feature-request template added under `.github/ISSUE_TEMPLATE/`. Help-section links point at them and at `https://discord.haze.sh`.

## Testing Decisions

- **Go**: new packages (`internal/logging`, `internal/presence/template`, `internal/startup`, the status bridge, the `Pause` flag, `Updater.LastSent`) held to the repo's existing table-driven test bar. Template engine: token substitution, unknown-token reporting, whitespace/separator collapse, each context's default reproducing the current string. Migration: a flat fixture maps to the expected v2 tree; a v2 fixture round-trips unchanged; a hand-broken file clamps to valid. Startup reconciler: tested against a fake registry writer, no real `HKCU` writes in CI. Status bridge: fake state manager and supervisor accessors, assert `GetStatus()` shape and that `status:changed` fires on change and not on no-op.
- **Frontend (Vitest + React Testing Library)**: template-preview rendering, config-form state and reducers, the global-default-vs-override resolution helper. No component snapshots, no E2E for v1.
- **Manual**: window hide/show, tray menu, start-with-Windows round trip, an actual update against a test release, first-run walkthrough with Discord/League toggled on and off.

## Out of Scope

- Beta / pre-release update channel. Stable only for v1.
- OS code-signing certificate. Builds are unsigned; SmartScreen warnings are accepted for now.
- macOS and Linux. Windows only.
- Auto-launching League itself. `AutoLaunchLeague` handling is not part of this milestone.
- Delta/patch updates. Full binary swap only.
- Per-queue (as opposed to per-`GameMode`) display config.
- Any presence-mapping behavior change beyond routing text through the template engine.
- MSIX / winget packaging.

## Further Notes

- Depends on [ADR-0004](../../docs/adr/0004-gui-process-hosts-the-daemon.md) (one process, no IPC) and [ADR-0005](../../docs/adr/0005-in-app-update-via-signed-binary-swap.md) (signed binary swap from GitHub Releases).
- `DEPENDENCIES.md`, the root `CLAUDE.md` frontend section, and `README.md` still describe shadcn/ui and an undecided frontend. Ticket 17 corrects them once the real stack is in place.
- Accent color and tray icon artwork are open. The accent is chosen during a `frontend-design` pass in ticket 11 for the user to review.
- The `GameMode` list for the per-mode UI comes from `pkg/types` at build time, not a hardcoded list here.
