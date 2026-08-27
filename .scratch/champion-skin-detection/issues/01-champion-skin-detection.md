# 01: Champion/skin detection for in-game presence

**What to build:** Real champion/skin/chroma art and an accurate elapsed-game timer on Discord presence during a match, replacing the current guaranteed-404 image bug, for every non-TFT `GameFlowInProgress` mode (Summoner's Rift, ARAM, URF, Brawl, League Classic, ARAM: Mayhem, Ultimate Spellbook, Arena, Swarm). TFT players see a correct level instead of a stuck 0. Full detail in `../spec.md`.

**Blocked by:** None (can start immediately)

**Status:** done

- [x] `internal/championdata` resolves a `rawChampionName`/`rawSkinName`/`skinID` triple into Data Dragon's `id` (URL-safe) and `name` (display), the resolved base skin name, and chroma name if applicable, using the three-tier raw-name fallback and the Meraki-then-DDragon chroma/base-skin resolution described in the spec. Champion JSON and the Data Dragon version string are cached in memory. HTTP transport is injected; unit tests use a fake transport, no real DDragon/Meraki access.
- [x] `internal/livegame` resolves the local player's riotId via `/liveclientdata/activeplayer`, matches them into `/liveclientdata/allgamedata`'s `allPlayers`, polls `/liveclientdata/playerscores` for KDA/CS, polls `/liveclientdata/gamestats` for elapsed `gameTime`, and independently polls `/liveclientdata/activeplayer` for TFT `level`. HTTP transport is injected; unit tests cover the reference implementation's documented edge cases (default-skin format, K/DA-variant format, loading-screen placeholder values) against fake fixtures, no real League client.
- [x] `daemon.go`'s `presenceLoop` starts the Live Client Data poller the instant it enters `GameFlowInProgress` and stops it the instant the phase changes, mirroring the existing placeholder-ticker start/stop.
- [x] For non-TFT modes: champion/skin/chroma resolve once per match (cached, not re-resolved every tick); KDA/CS/timer poll every `StatsPollingInterval` tick; while champion data is unresolved (loading screen) or on a transient poll failure, the poller skips writing `State` that tick rather than zeroing/clearing presence, relying on `state.Manager.Update`'s existing "unchanged, skip" behavior.
- [x] For TFT: an independent per-tick poll of `activeplayer.level` feeds `state.Manager.UpdateInGameStats`, with no riotId/allgamedata/championdata involvement.
- [x] `State.ChampionName` is split into a URL-safe ID field (feeds `GetChampionSkinURL`) and a display-name field (feeds `FormatSkinName`), sourced from Data Dragon's own `id`/`name` fields — no hand-maintained conversion map.
- [x] `BuildInGamePresence`'s `Start` field uses the poller-supplied timestamp computed from real `gameTime`, not `time.Now()`, so the visible timer no longer resets every heartbeat.
- [x] `internal/discord/assets.go`'s `animatedSkins` map includes the 5 entries missing relative to the reference implementation: `Jinx_60`, `Kaisa_71`, `Mordekaiser_54`, `Morgana_80`, `Sett_66`.
- [x] Outbound HTTP calls in both new packages use short timeouts (a few seconds), not a blanket 15s.
- [x] Locale stays hardcoded (`en_US` / `en-US`); no new `Config` field.
- [x] No test in either new package or the poller-lifecycle wiring depends on a real League client, real Discord client, or live internet access.
