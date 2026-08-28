# 19: GameModes() can silently drift from the GameMode const block

**What to fix:** `pkg/types/queue.go` declares the `GameMode` constants in one block and `GameModes()` returns a hand-written slice of them. The two are not linked. Add a new mode const without appending it to `GameModes()` and:

- `ValidGameMode("NEWMODE")` returns false
- `resolveDisplay` (`internal/discord/mapper.go`) treats it as unknown and returns `Display.Default`
- any per-mode override the user configured for that mode is silently ignored, with nothing logged or failed

**Blocked by:** none

**Status:** done

- [x] `GameModes()` and `ValidGameMode` cannot disagree with the set of declared `GameMode` constants
- [x] Approach is either a single declaration both derive from, or a test that fails when a const is added without updating the list
- [x] Existing callers of `GameModes()` / `ValidGameMode` unaffected

## Comments

Took the guard-test route rather than a refactor: `GameModes()` stays a plain slice literal, and `TestGameModes_MatchesConstBlock` parses `queue.go` with `go/ast`, collects every `GameMode`-typed string constant, and asserts the set matches `GameModes()` both ways (missing entry and stray entry). Verified it fails when a const is dropped from the list. No production code changed.
