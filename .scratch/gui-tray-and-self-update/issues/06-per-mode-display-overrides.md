# 06: Per-mode display overrides

**What to build:** A resolution helper that, given a `GameMode` and the config, returns the effective `ShowRank` and `ShowStats`: the `Display.Modes[mode]` override if present, otherwise `Display.Default`. Wire it into `discord/mapper.go` so presence honors per-mode settings. The mode key is the `GameMode` display concept, never a queue ID. No other display setting is per-mode.

**Blocked by:** 02

**Status:** done

- [x] A pure helper `resolveDisplay(cfg, mode) (showRank, showStats bool)` with default-vs-override logic
- [x] `discord/mapper.go` calls it instead of reading global `ShowRank`/`ShowStats` directly
- [x] A mode with no `Display.Modes` entry inherits `Display.Default` exactly
- [x] The set of valid mode keys is derived from `pkg/types`, not hardcoded in this package
- [x] Tests: default inheritance, override wins, unknown mode falls back to default, mapper output reflects the resolved values

## Comments

### Implementation notes

- **`resolveDisplay` lives in `discord/mapper.go`**, returns the `(showRank, showStats)` pair the ticket names. It reuses `config.DisplayConfig.Resolve` for the default-vs-override merge and only adds the known-mode gate: an unknown `GameMode` (or one with a stale override stored under a bogus key) returns `Display.Default` untouched.
- **Builders were left alone.** All seven per-mode reads live in `discord/presence.go` as `cfg.Display.Default.ShowRank/ShowStats`. Rather than thread a resolved struct through seven builder signatures and their tests, `MapStateToPresence` calls `withResolvedDisplay(st.GameMode, cfg)` once and hands the builders a shallow copy whose `Display.Default` already holds the resolved values. The copy shares the `Modes`/`Templates` maps by reference; nothing mutates them.
- **New in `pkg/types`:** `GameModes() []GameMode` and `ValidGameMode(GameMode) bool`. `GameModes()` is the single source both this resolver and the future per-mode settings UI (ticket 12) read; nothing hardcodes a mode list in `discord`.
- **Mode key is the `GameMode` string** (`"CLASSIC"`, `"CHERRY"`, ...), never a queue ID. The existing `config.DisplayConfig_Resolve` test uses `"ARENA"` as a generic key; that is fine for exercising the merge in isolation but is not a real mode key, so `resolveDisplay` would treat it as unknown.
