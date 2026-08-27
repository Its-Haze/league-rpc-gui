# 15: Frontend: Home dashboard and first-run onboarding

**What to build:** The Home section. A status dashboard: the three connection states with plain-language explanations, the current `GameFlowPhase`, and a live preview card rendering the presence the `Updater` last sent (from `GetStatus()`), updating on `status:changed`. A "Test presence" button calling `App.TestPresence()`. A first-run inline walkthrough (not a modal): what the app does, a step that shows live checkmarks as Discord and League connect, a step to pick what to show, a step offering start-with-Windows. Dismissible and re-openable from Help.

**Blocked by:** 11, 08

**Status:** ready-for-agent

- [ ] Dashboard shows the three connection states with short explanations and the current phase
- [ ] Live preview card mirrors `GetStatus()` last-sent presence and updates on `status:changed`
- [ ] "Test presence" button calls `App.TestPresence()` and indicates the ~30s test window
- [ ] First-run walkthrough renders inline on first launch only, with live Discord/League checkmarks
- [ ] Walkthrough is dismissible; a Help entry reopens it
- [ ] A flag (config or a small state file) records that onboarding was completed
- [ ] Vitest: onboarding step progression, preview binding wiring
