# Champion/skin detection: in-game presence detail

Status: ready-for-agent

## Problem Statement

`BuildInGamePresence` (`internal/discord/presence.go`) already calls `GetChampionSkinURL(st.ChampionName, st.SkinID)` and formats `st.ChampionName`/`st.SkinName`/`st.ChromaName` for display, and `state.State` already has fields for all of it. But nothing populates those fields during `GameFlowInProgress`: no code polls Live Client Data, no code resolves a champion/skin/chroma name from Data Dragon or Meraki. The result is a guaranteed-404 image URL (`.../tiles/_0.jpg`), which Discord renders as a broken "?" tile instead of the champion's skin art, for every player, in every match. `BuildInGamePresence` also stamps `Start: time.Now().Unix()` fresh on every call, which resets Discord's visible elapsed-game timer back to `0:00` on every heartbeat/poll cycle instead of counting up from the real game start.

This is the explicit follow-on to the daemon-skeleton spec, which proved the Daemon's core lifecycle (connects, reconnects, shows correct presence for every phase except in-game detail) while deliberately excluding `internal/livegame` (Live Client Data) and `internal/championdata` (Data Dragon/Community Dragon) as out of scope. This spec builds those two packages and wires them into the existing `GameFlowInProgress` path.

## Solution

Add two new packages and a poller that the Daemon starts only while `GameFlowInProgress` is active:

- **`internal/livegame`**: talks to the unauthenticated Live Client Data API at `127.0.0.1:2999`. Identifies the local player via `/liveclientdata/activeplayer`'s `riotId`, matches that riotId into `/liveclientdata/allgamedata`'s `allPlayers` list to read `rawChampionName`/`rawSkinName`/`skinID`, and polls `/liveclientdata/playerscores?riotId=` for kills/deaths/assists/creepScore. Computes an accurate presence `Start` from `/liveclientdata/gamestats`'s `gameTime`. For TFT, bypasses all of the above and polls only `/liveclientdata/activeplayer`'s `level`.
- **`internal/championdata`**: resolves a raw champion identifier (from `rawChampionName`/`rawSkinName`, using the same three-tier fallback the reference implementation uses for the API's inconsistent field formats) into a champion's Data Dragon `id` (URL-safe, feeds asset URLs) and `name` (display text), plus the resolved base skin name and, if applicable, chroma name (cross-referencing Community Dragon's Meraki Analytics champion JSON against Data Dragon's skin list, the same way the reference implementation's `find_base_skin_and_chroma` does). Caches champion JSON and the Data Dragon version string in memory.
- Both packages take an injected HTTP-client-like interface (matching the `ProcessChecker`/`Connector` pattern already used in `internal/daemon`), so both are unit-testable against fake responses with no real League client, DDragon, or Meraki access.

The Daemon starts this poller the moment `presenceLoop` (in `daemon.go`) sees `GameFlowInProgress`, the same way it already starts/stops the placeholder ticker, and stops it the moment the phase changes. Poll results write into `state.Manager` via the existing `UpdateChampion`/`UpdateInGameStats` methods; `MapStateToPresence`/`BuildInGamePresence`/`BuildTFTInGamePresence` are not restructured, only fed real data for the first time.

## User Stories

1. As a player, I want my Discord presence to show my actual champion and skin art during a match, so that friends see what I'm playing instead of a broken image icon.
2. As a player playing a chroma, I want my presence to show the chroma's name alongside the base skin, so that my presence accurately reflects what I selected.
3. As a player, I want my in-game Discord presence timer to count up from when the match actually started, so that it doesn't visibly reset every time presence resends.
4. As a player during the loading screen (before champion data is available), I want my presence to keep showing champ-select's last state rather than a broken or blank card, so that I never see a jarring flash of missing data.
5. As an Arena or Swarm player, I want my champion/skin to show correctly too, so that the same broken-image bug doesn't persist for those modes just because they're not Summoner's Rift.
6. As a TFT player, I want my presence to show my actual level, so that the "In Game · lvl: N" text isn't stuck at 0.
7. As a developer, I want champion/skin/chroma resolution and Live Client Data polling to be fully testable in CI without a real League client or live internet access, so that this logic doesn't rot untested.
8. As a developer, I want the small set of animated Ultimate skins we already special-case to be complete, so that all animated skins the reference implementation knows about still animate here.

## Implementation Decisions

- **New package `internal/livegame`**: owns all `127.0.0.1:2999` HTTP calls (self-signed cert, `InsecureSkipVerify: true`, same as the LCU client). Exposes whatever the poller needs to resolve one match's data; the HTTP transport is injected (an interface satisfied by `*http.Client` in production, a fake in tests) per the existing `ProcessChecker`/`Connector` injection pattern in `internal/daemon`.
- **New package `internal/championdata`**: owns all Data Dragon (`ddragon.leagueoflegends.com`) and Community Dragon/Meraki (`cdn.merakianalytics.com`) HTTP calls, plain HTTP (no League process needed). Same injected-transport pattern as `internal/livegame`.
  - Champion name resolution ports the reference implementation's three-tier fallback: `rawChampionName.split("_")[-1]` (primary), then `rawSkinName.split("_")[-1]` (default-skin fallback), then `rawSkinName.split("_")[-2]` (K/DA-variant-skin fallback), each tried against Data Dragon until one resolves.
  - `State.ChampionName` is split into two concepts at the resolution layer: Data Dragon's `id` field (URL-safe, e.g. `"Chogath"`) feeds `GetChampionSkinURL`, and its `name` field (display text, e.g. `"Cho'Gath"`) feeds `FormatSkinName`. This replaces the reference implementation's hand-maintained `CHAMPION_NAME_CONVERT_MAP` (a ~60-entry table duplicating data Data Dragon's own response already carries) with the API's own display name.
  - Skin/chroma resolution: given `skinID`, first check Meraki's chroma list (matching `chroma.id % 1000 == skinID`) for a chroma name and its parent base-skin number; if not found there, walk Data Dragon's skin list for the highest `num <= skinID` whose name has no parenthetical suffix (its own base-skin heuristic). The resulting skin tile URL is built directly from Data Dragon's own skin list (`{ddragon}/img/champion/tiles/{champUrlID}_{baseSkinNum}.jpg`) — no live HTTP HEAD verification is performed; the DDragon response already tells us which skin numbers are valid.
  - Champion JSON is cached in memory, keyed by champion ID + Data Dragon version. The Data Dragon version string itself is also cached (fetched once per daemon run, not re-fetched on every resolution attempt) rather than hit on every lookup.
  - Locale is hardcoded to `en_US` (Meraki: `en-US`) as a literal, not a `Config` field. `CONTEXT.md`'s existing Locale entry ("auto-detected from the OS process, never a user-facing setting") stands unchanged; real locale threading is deferred to its own spec.
  - Animated-skin table (`internal/discord/assets.go`'s `animatedSkins`) gets the 5 entries missing relative to the reference implementation's `ANIMATED_SKINS`: `Jinx_60`, `Kaisa_71`, `Mordekaiser_54`, `Morgana_80`, `Sett_66`.
- **Poller lifecycle**: `daemon.go`'s `presenceLoop` starts an `internal/livegame` poller goroutine the instant it enters `modeConnected` with `GameFlowInProgress`, and stops it the instant the phase changes, mirroring the existing placeholder-ticker start/stop. The poller runs on `config.StatsPollingInterval`.
  - **Non-TFT path**: on the first tick, resolve riotId via `/liveclientdata/activeplayer`, then match it into `allgamedata`'s `allPlayers` and resolve champion/skin/chroma via `internal/championdata`. Once resolved, cache the result for the rest of the match (champion/skin/chroma don't change mid-game) and stop re-resolving. Every tick (including before resolution succeeds), poll `playerscores?riotId=` (reusing the cached riotId) for KDA/CS, and `/liveclientdata/gamestats` for `gameTime` to compute an accurate `Start`.
  - **TFT path**: independent, no riotId/allgamedata/championdata involvement — each tick polls `/liveclientdata/activeplayer` directly for `level` and calls `state.Manager.UpdateInGameStats`.
  - **Unresolved champion (loading screen, placeholder API values like `"Name"`/`"Unknown"`/`""`)**: the poller simply does not call `UpdateChampion` until resolution succeeds. `state.Manager.Update`'s existing "unchanged, skip" semantics mean the prior presence (most likely champ-select's) stays on screen with no new logic required in `presence.go`/`mapper.go`.
  - **Transient tick failure** (network hiccup, brief error) after champion/skin/chroma have already resolved: skip writing `State` that tick rather than zeroing KDA/CS/timer or clearing presence, same "leave last known state" posture the daemon-skeleton spec already established for League-process detection.
  - Outbound HTTP calls (Live Client Data, Data Dragon, Meraki) use short timeouts (a few seconds, not the reference implementation's blanket 15s), since a hang inside a poll tick would otherwise stall the shared poll loop; a timeout is just "try again next tick."
- **Game-mode routing**: every non-TFT `GameFlowInProgress` game (Summoner's Rift, ARAM, URF, Brawl, League Classic, ARAM: Mayhem, Ultimate Spellbook, Arena, Swarm) resolves champion/skin/chroma and routes through the existing `BuildInGamePresence` (KDA+CS layout). Arena's and Swarm's own distinct level+gold stat lines (and Arena's rank-swap) are not built in this spec — see Out of Scope.
- `MapStateToPresence`, `BuildInGamePresence`, `BuildTFTInGamePresence`, and `ShouldClearPresence` are not restructured; `BuildInGamePresence`'s `Start` field is changed to use the poller-supplied accurate timestamp instead of `time.Now().Unix()`.

## Testing Decisions

- **`internal/livegame`**: unit tests inject a fake HTTP transport serving canned Live Client Data JSON (including the reference implementation's documented edge cases: default-skin format, K/DA-variant format, placeholder loading-screen values). Assert riotId resolution, player matching, KDA/CS extraction, TFT level extraction, and that a failed/incomplete response results in no `State` write rather than a zeroed one.
- **`internal/championdata`**: unit tests inject a fake HTTP transport serving canned Data Dragon/Meraki JSON. Assert the three-tier raw-name fallback chain, the chroma-vs-base-skin resolution logic (including the Meraki-first-then-DDragon-heuristic ordering, since DDragon alone can misidentify a chroma as a base skin), skin URL construction, and that champion JSON/version caching avoids repeat fetches within a session.
- **Poller lifecycle**: extend the existing `daemon` test seam (`Daemon.Run`, fake connectors/clock) to assert the Live Client Data poller starts only on `GameFlowInProgress` and stops on any other phase.
- No test relies on a real League client, real Discord client, or live internet access, per the daemon-skeleton spec's existing CI-independence requirement.

## Out of Scope

- Spectator mode (`Watching` phase) — not even defined in `pkg/types/gameflow.go` yet.
- Arena's and Swarm's dedicated stat-line layouts (level, gold, Arena's rank-swap) — they get champion/skin resolution and the generic KDA+CS layout only, in this spec.
- Locale auto-detection and any config surface for it — `en_US`/`en-US` stays hardcoded; `CONTEXT.md`'s existing Locale glossary entry is unchanged.
- League Classic and ARAM: Mayhem display-name/icon work beyond what already exists in `internal/discord/assets.go`.
- Any Wails GUI or frontend work.

## Further Notes

- Builds directly on the daemon-skeleton spec's `Daemon`/`presenceLoop`/`state.Manager` machinery; no changes to the Connection Supervisors, `Updater`, or the non-in-game presence builders.
- The HTTP-client-injection pattern for `internal/livegame`/`internal/championdata` follows the same shape as `ProcessChecker`/`Connector` in `internal/daemon` (see `internal/daemon/lcu_supervisor.go`), keeping the whole codebase consistent on "inject the seam, don't touch a real external system in tests."
- Locale, spectator mode, and Arena/Swarm's own stat-line parity remain explicitly reserved for future specs, as the daemon-skeleton spec's Further Notes already anticipated.
