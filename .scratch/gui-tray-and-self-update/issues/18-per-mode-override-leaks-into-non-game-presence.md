# 18: Per-mode display override leaks into idle and post-game presence

**What to fix:** `MapStateToPresence` (`internal/discord/mapper.go:19`) calls `withResolvedDisplay(st.GameMode, cfg)` for every gameflow phase, so a `Display.Modes[mode]` override governs presence even when the user is not in that mode's game.

`state.NewState()` sets `GameMode` to `GameModeClassic` and nothing resets it: `resetInGameData` in `internal/state/manager.go` leaves `GameMode` alone. So:

- On a fresh launch while idle in client, `GameMode` is still `CLASSIC`. If the user set `Display.Modes["CLASSIC"] = {ShowRank:false}`, `BuildInClientPresence` hides the rank on the idle card.
- After any match, `GameMode` retains the last mode played (e.g. `ARAM`). Its override then applies to the post-game, lobby, and queue presence until the next game sets a new mode.

Before ticket 6 the builders read `Display.Default` directly, so the idle card was never affected by a per-mode override.

**Blocked by:** none

**Status:** done

- [x] Per-mode resolution only applies to phases where the mode is actually being played (in-game, spectating, and arguably champ select / lobby for that mode); idle and post-game use `Display.Default`
- [x] Decide and document which phases count as "in this mode" (the `GameMode` field is only meaningful once a game or lobby for that mode exists)
- [x] Test: idle-in-client presence ignores a `Display.Modes["CLASSIC"]` override
- [x] Test: post-game presence (`GameFlowEndOfGame`, `GameFlowWaitingForStats`) ignores the last-played mode's override
- [x] Test: in-game presence still honors the override

## Comments

Impact was latent, not observed: no in-client / post-game builder reads `ShowRank` or `ShowStats` today, so the stale-mode `withResolvedDisplay` call was harmless by accident. The fix makes it correct by design.

`MapStateToPresence` now gates the per-mode resolve on `perModeDisplayApplies(phase)`, a pure helper that returns `!phase.IsInClient()`. `GameFlowPhase.IsInClient()` is exactly the idle + post-game set (`None`, `WaitingForStats`, `PreEndOfGame`, `EndOfGame`); every phase that renders rank/stats (Lobby, queue, champ select, in-game, spectating) has a live `GameMode` from the lobby config or live game. Tests: `TestPerModeDisplayApplies` (phase table) and `TestMapStateToPresence_PostGameIgnoresStaleModeOverride`.
