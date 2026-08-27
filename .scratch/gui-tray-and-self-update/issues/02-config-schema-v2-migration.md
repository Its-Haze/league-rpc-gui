# 02: Config schema v2 and one-time migration

**What to build:** Move `config.Config` from the flat struct to a nested tree with a top-level `SchemaVersion int`. Groups: `Display` (`Default.{ShowRank,ShowStats}`, `Modes map[string]ModeOverride`), `Presence` (`ShowEmojis`, `ShowInClient`, `Idle`, `Templates map[string]TemplatePair`), `Behavior` (`LaunchAtStartup`), `Advanced` (`UpdateInterval`, `StatsPollingInterval`, `DebugMode`), plus top-level `DiscordAppID` and `Theme`. `Load` detects a versionless file, migrates it into the tree, and writes it back. Every current consumer of `config.Config` moves to the new paths.

**Blocked by:** None (can start immediately)

**Status:** done

- [x] `Config` is the nested tree with `SchemaVersion` (`Display` / `Presence` / `Behavior` / `Advanced` groups, top-level `DiscordAppID` / `Theme`); `DefaultConfig()` returns a fully-populated tree
- [x] `Load` migrates a flat (versionless) file via `migrate.go`: old fields map to their new locations, `SchemaVersion` is stamped, the upgraded file is written back (best-effort `Save`)
- [x] A v2 file round-trips through `Load`/`Save` unchanged (`TestLoad_V2FileRoundTrips`, JSON compare)
- [x] `clamp()` has a branch per numeric/enum field (`DiscordAppID`, `Theme`, both intervals) plus nil-map repair
- [x] `Validate()` covers `Theme` and the moved numeric fields, reports every problem at once, proven non-mutating by a marshal-compare test
- [x] `Theme` is `system` / `light` / `dark`, defaulting to `system`
- [x] `discord/mapper.go`, `discord/presence.go`, `discord/updater.go`, `lcu/client.go`, `livegame/poller.go`, `cmd/league-rpc` moved to the new paths; `AutoLaunchLeague` / `LeaguePath` kept under `Behavior`
- [x] Tests: flat fixture to expected tree, v2 round-trip, broken file clamps to valid, `Validate` table, `Resolve` override helper
