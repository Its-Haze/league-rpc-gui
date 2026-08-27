# 13: Frontend: Behavior screen

**What to build:** The Behavior section. A "start with Windows" toggle bound to `Behavior.LaunchAtStartup` that reconciles the registry immediately (ticket 07). A short explainer of the tray behavior (X hides, right-click to quit). In-client / idle presence controls (`Presence.ShowInClient` and the `Idle` settings). Update settings: the check cadence is fixed, so this is where the manual "Check for updates" button and the current channel (stable) live.

**Blocked by:** 11, 07

**Status:** ready-for-agent

- [ ] "Start with Windows" toggle writes `Behavior.LaunchAtStartup` and triggers an immediate registry reconcile
- [ ] Tray-behavior explainer text is present and accurate
- [ ] `Presence.ShowInClient` and `Idle` controls bound to the config tree
- [ ] "Check for updates" button calls the ticket-09 manual-check binding and shows the result
- [ ] Channel shown as stable, with no beta option for v1
- [ ] Vitest: toggle state and the reconcile call with a mocked binding
