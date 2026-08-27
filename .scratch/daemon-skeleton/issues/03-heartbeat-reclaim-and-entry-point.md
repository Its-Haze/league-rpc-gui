# 03: Heartbeat/reclaim + the real entry point

**What to build:** Extend the existing `discord.Updater` per ADR-0001: its one-shot debounce becomes a repeating ticker that keeps resending the current presence (heartbeat) after the first real send, with a faster burst right after each real change (reclaim), so League's own native Discord integration can't silently overwrite ours. Then wire `cmd/league-rpc/main.go` with real `discord.Client`/`lcu.Client` instances calling `Daemon.Run`, listening for OS shutdown signals as the one path by which `Run` is expected to return. This is the point the whole feature becomes something you can actually `go run`.

**Blocked by:** 02

**Status:** ready-for-agent

- [ ] `discord.Updater`'s one-shot debounce timer becomes a repeating ticker that continues resending the current presence (heartbeat) after the first real send
- [ ] After each real presence change, `Updater` fires a short burst of faster resends (reclaim) before settling back to the normal heartbeat cadence
- [ ] All presence sends (heartbeat, reclaim, placeholder rotation, and real updates alike) go through `Updater`'s existing mutex; no new lock is introduced, and `discord.Client.UpdatePresence`/`ClearPresence` are only ever called from within `Updater`
- [ ] Heartbeat/reclaim timing is tested via an injected clock/ticker abstraction, not real `time.Sleep`/wall-clock waits
- [ ] `cmd/league-rpc/main.go` constructs a real `Daemon` with real `discord.Client` and `lcu.Client`, and calls `Daemon.Run`
- [ ] `main.go` listens for OS shutdown signals (Ctrl+C/SIGTERM) and cancels `Daemon.Run`'s context on receipt. This is the only path by which `Run` is expected to return
- [ ] Running `go run ./cmd/league-rpc` with real Discord and League running connects successfully, shows presence, and shuts down cleanly on Ctrl+C
