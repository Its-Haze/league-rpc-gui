# 06: Per-mode display overrides

**What to build:** A resolution helper that, given a `GameMode` and the config, returns the effective `ShowRank` and `ShowStats`: the `Display.Modes[mode]` override if present, otherwise `Display.Default`. Wire it into `discord/mapper.go` so presence honors per-mode settings. The mode key is the `GameMode` display concept, never a queue ID. No other display setting is per-mode.

**Blocked by:** 02

**Status:** ready-for-agent

- [ ] A pure helper `resolveDisplay(cfg, mode) (showRank, showStats bool)` with default-vs-override logic
- [ ] `discord/mapper.go` calls it instead of reading global `ShowRank`/`ShowStats` directly
- [ ] A mode with no `Display.Modes` entry inherits `Display.Default` exactly
- [ ] The set of valid mode keys is derived from `pkg/types`, not hardcoded in this package
- [ ] Tests: default inheritance, override wins, unknown mode falls back to default, mapper output reflects the resolved values
