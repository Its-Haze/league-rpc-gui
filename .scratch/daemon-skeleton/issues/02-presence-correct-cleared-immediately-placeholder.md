# 02: Presence: correct, cleared immediately, placeholder while launching

**What to build:** Wire `Daemon.Run` to the existing, unmodified `discord.MapStateToPresence`/`discord.ShouldClearPresence` and `internal/state`, via the existing `Updater`. Once both supervisors report connected, the correct presence appears for every in-scope `GameFlowPhase` (client-idle, lobby, queue, champ select). The moment the LCU disconnects with no League process detected, presence clears immediately, no grace period, no fallback idle card. While League's process is detected running but the LCU hasn't connected yet, a rotating "Launching League..." placeholder shows instead of nothing, stopping the instant real presence takes over. This is the first ticket that's actually demoable end to end.

**Blocked by:** 01

**Status:** ready-for-agent

- [ ] Once both Connection Supervisors report connected, `Daemon.Run` drives presence through the existing `MapStateToPresence`/`ShouldClearPresence` and `Updater`: correct presence appears for InClient, Lobby, Matchmaking/ReadyCheck (Queue), and ChampSelect/GameStart
- [ ] The moment the LCU Connection Supervisor reports disconnected and no League process is detected, `Updater` clears presence immediately, no grace period, no fallback idle card
- [ ] While the League process is detected running but the LCU hasn't connected yet, a rotating "Launching League..." placeholder presence is shown instead of nothing
- [ ] The placeholder stops the instant the LCU connects, and real presence takes over from there
- [ ] `MapStateToPresence`/`ShouldClearPresence` themselves are unchanged by this ticket. Their existing signatures and pure-function shape are preserved; only how/when `Daemon` calls them is new
- [ ] New Daemon-level behavior (clear-on-disconnect, placeholder trigger, transition to real presence) is tested through the `Daemon.Run` seam with fakes, not real Discord/LCU clients
