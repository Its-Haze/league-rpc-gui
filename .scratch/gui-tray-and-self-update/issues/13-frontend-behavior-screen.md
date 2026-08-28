# 13: Frontend: Behavior screen

**What to build:** The Behavior section. A "start with Windows" toggle bound to `Behavior.LaunchAtStartup` that reconciles the registry immediately (ticket 07). A short explainer of the tray behavior (X hides, right-click to quit). In-client / idle presence controls (`Presence.ShowInClient` and the `Idle` settings). Update settings: the check cadence is fixed, so this is where the manual "Check for updates" button and the current channel (stable) live.

**Blocked by:** 11, 07

**Status:** done

- [x] "Start with Windows" toggle writes `Behavior.LaunchAtStartup` and triggers an immediate registry reconcile
- [x] Tray-behavior explainer text is present and accurate
- [x] `Presence.ShowInClient` and `Idle` controls bound to the config tree
- [x] "Check for updates" button calls the ticket-09 manual-check binding and shows the result
- [x] Channel shown as stable, with no beta option for v1
- [x] Vitest: toggle state and the reconcile call with a mocked binding

**Notes:** The "immediate registry reconcile" is the daemon's own
`reconcileStartupOnChange` goroutine (`cmd/league-rpc-gui/main.go`), which was already watching
`config.Store` for `Behavior.LaunchAtStartup` flips since ticket 07; the toggle here just calls
`ApplySettings` like every other field and the reconcile follows automatically, so no new
binding was needed. `withLaunchAtStartup`/`withShowInClient`/`withIdleText` in
`lib/behaviorPatch.ts` are the patch builders Vitest exercises, kept separate from the component
so the "reconcile call" (really: the patch shape `ApplySettings` receives) is testable without
mounting React or mocking the Wails runtime.
