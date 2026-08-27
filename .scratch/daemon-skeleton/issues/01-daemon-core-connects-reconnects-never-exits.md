# 01: Daemon core: connects, reconnects, never exits

**What to build:** The Daemon's core lifecycle, with no presence output yet. `internal/process` detects whether Discord or League's own process is running. A Discord Connection Supervisor and an LCU Connection Supervisor each retry their respective connection forever, independently of each other, never exiting on failure. `Daemon.Run(ctx)` is the single top-level seam: it starts both supervisors and blocks until `ctx` is canceled, disconnecting both cleanly on the way out. This proves the "never needs a restart" contract in isolation, before any presence behavior is layered on top.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] `internal/process` exposes a check for whether any process in a given name list is currently running, using `gopsutil`, testable against a mocked process-listing call (no real Discord/League process required in CI)
- [ ] The Discord Connection Supervisor waits for Discord's process (via `internal/process`) before each `discord.Client.Connect()` attempt, retries indefinitely on a short fixed interval, and never returns/exits on failure
- [ ] The LCU Connection Supervisor wraps `lcu.Client.Connect()`/`Disconnect()` in a retry-forever loop (`Connect()` already blocks via `AwaitConnection: true`), and exposes "League process detected" as an observable signal for later use (ticket 02's placeholder trigger)
- [ ] The Discord and LCU Connection Supervisors run independently: one failing repeatedly never blocks or delays the other from succeeding (covered by a test)
- [ ] `Daemon.Run(ctx)` starts both supervisors and blocks until `ctx` is canceled; it never returns on its own due to a connection error (this is the single seam for the whole daemon-skeleton effort, see spec's Further Notes)
- [ ] On `ctx` cancellation, `Daemon.Run` stops both supervisors and disconnects both clients, best-effort (matching the existing "never let cleanup panic on the way out" posture)
- [ ] All new behavior is tested via injected fakes (fake Discord/LCU connectors, mocked process-listing); no test depends on a real Discord or League process
