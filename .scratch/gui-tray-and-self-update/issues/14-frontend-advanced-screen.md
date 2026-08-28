# 14: Frontend: Advanced screen

**What to build:** The Advanced section. Discord Application ID: a preset picker (from `GetPresets`) plus a custom-entry field. RPC update interval and stats polling interval, with the min/max from `config` enforced in the UI. A debug-logging toggle bound to `Advanced.DebugMode`.

**Blocked by:** 11, 02

**Status:** done

- [x] Discord App ID: preset dropdown from `GetPresets` and a validated custom field; empty is rejected inline
- [x] Update-interval and stats-polling-interval inputs clamped to the `config` bounds, with the bound shown
- [x] `Advanced.DebugMode` toggle; changing it takes effect without a restart (log level follows it)
- [x] All changes go through `ApplySettings` and surface `Validate()` errors inline
- [x] Vitest: bounds enforcement, preset-vs-custom switching

**Notes:** Bounds (`lib/advancedBounds.ts`) originally mirrored `config.MinUpdateInterval` /
`MaxUpdateInterval` / `MinStatsPollingInterval` / `MaxStatsPollingInterval` as hand-kept
constants; a `/code-review` pass on the drift risk led to adding `App.GetConfigBounds()` (mirrors
`GetGameModes()`/`GetTemplateTokens()`), so the screen now reads bounds from the backend and the
hand-kept constants only serve as the fallback shown before that first resolves. The preset/custom
switch is a `Select` plus a custom-entry field that only appears once the id doesn't match any
preset (`presetSelectValue`). Debug mode's "no restart" requirement needed a backend fix: the
logger was setting its level once at construction (`.Level(level)`), which can't change
afterward. Switched `internal/logging.New` to leave the per-logger level unset and drive
`zerolog.SetGlobalLevel` instead (`logging.SetDebug`), and added a `reconcileDebugOnChange`
watcher (later generalized into `watchConfigField[T]`) in `cmd/league-rpc-gui/main.go` that calls
it on every `Advanced.DebugMode` flip. The toggle itself just calls `ApplySettings` like
everything else here.

**Follow-up:** added a "How do I get a custom Application ID?" walkthrough (`<details>`
disclosure linking to the Discord Developer Portal) below the custom-entry field, plus a
best-effort name resolution: `internal/discordapp` (new package, mirrors `internal/updates`'s
injectable-`HTTPDoer` pattern) calls Discord's public `applications/{id}/rpc` endpoint and the
screen shows "Resolves to <name>" once it's typed. Wired through `App.GetApplicationName` /
`WithAppNameLookup`, cached by id on the frontend (`hooks/useDiscordAppName`) so Home (ticket 15)
and this screen don't each refetch the same id. A failed lookup (bad id, Discord unreachable) is
silent, not an error state, since this is a nice-to-have, not load-bearing.
