# 03: Logging to a rotating file and an in-memory ring buffer

**What to build:** A new `internal/logging` package that builds a zerolog logger writing to both a `lumberjack`-rotated file at `%APPDATA%\league-rpc\logs\league-rpc.log` (size-capped, keep 2 old files) and an in-memory ring buffer holding the last ~500 lines. Level follows `Advanced.DebugMode`. Expose a way to read the buffer and to subscribe to new lines. The GUI entry point uses this logger; `cmd/league-rpc` keeps its console writer.

**Blocked by:** None (the package); GUI wiring of `log:line` events lands with ticket 16

**Status:** done

- [x] `internal/logging.New` returns a `*Sink` whose `Logger` is a `zerolog.Logger` fanned out via `MultiLevelWriter` to the file and the ring
- [x] File rotation via `lumberjack`: `MaxSize` cap, `MaxBackups: 2`, at `%APPDATA%\league-rpc\logs\league-rpc.log`, dir created by `LogDir()`
- [x] `Ring` holds ~500 lines (`DefaultRingSize`), mutex-guarded read/append, wraps and drops oldest on overflow
- [x] `Ring.RecentLines()` returns oldest-first; `Ring.Subscribe()` returns a channel of new lines (non-blocking sends)
- [x] Level is `DebugLevel` when `Options.Debug` (fed from `Advanced.DebugMode` in `cmd/league-rpc-gui`), else `InfoLevel`
- [x] `cmd/league-rpc-gui/main.go` uses `logging.New`; `cmd/league-rpc` keeps its `zerolog.ConsoleWriter`
- [x] Tests: `TestRing_WrapsAtCapacity`, `TestRing_ConcurrentAppendAndReadIsRaceFree`, `TestNew_LevelFollowsDebugFlag`, file+ring fan-out. Note: `-race` needs cgo/gcc, absent here; the concurrency test still runs without it
