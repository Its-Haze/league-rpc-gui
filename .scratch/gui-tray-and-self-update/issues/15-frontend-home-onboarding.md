# 15: Frontend: Home dashboard and first-run onboarding

**What to build:** The Home section. A status dashboard: the three connection states with plain-language explanations, the current `GameFlowPhase`, and a live preview card rendering the presence the `Updater` last sent (from `GetStatus()`), updating on `status:changed`. A "Test presence" button calling `App.TestPresence()`. A first-run inline walkthrough (not a modal): what the app does, a step that shows live checkmarks as Discord and League connect, a step to pick what to show, a step offering start-with-Windows. Dismissible and re-openable from Help.

**Blocked by:** 11, 08

**Status:** done

- [x] Dashboard shows the three connection states with short explanations and the current phase
- [x] Live preview card mirrors `GetStatus()` last-sent presence and updates on `status:changed`
- [x] ~~"Test presence" button calls `App.TestPresence()` and indicates the ~30s test window~~ removed at the user's request; see notes
- [x] First-run walkthrough renders inline on first launch only, with live Discord/League checkmarks
- [x] Walkthrough is dismissible; a Help entry reopens it
- [x] A flag (config or a small state file) records that onboarding was completed
- [x] Vitest: onboarding step progression, preview binding wiring

**Notes:** The "flag" is `Config.OnboardingComplete` (new field, `internal/config`), defaulting
false on `DefaultConfig()` and backfilled to true by `migrateFromV1` so an upgrading install
never sees the walkthrough it already lived through. The walkthrough (`OnboardingWalkthrough`)
is a plain `useReducer` over `lib/onboarding.ts`'s step list (welcome/connect/display/startup);
dismissing or finishing sets `onboarding_complete: true` through the same `useSettings().applyPatch`
every other screen uses. The Help screen's "Replay walkthrough" button just clears the flag again;
it does not also navigate to Home, so the user has to click over there to see it restart.
`hooks/useStatus.ts` is a new shared hook extracted from what `TopStrip` was already doing inline
(`GetStatus` + `status:changed`); `TopStrip` was refactored onto it too, to avoid a second parallel
implementation of the same subscribe logic.

**Follow-up:** the user asked to remove "Test presence" entirely as unnecessary. Removed
end-to-end rather than just hiding the button: `App.TestPresence`, the `TestPresenter` interface,
`Daemon.TestPresence` and its `testActive`/`testDuration`/`testSignal` fields (which had their own
carve-out in `presenceLoop`), `Updater.PushSample`, and `discord.BuildTestPresence`, plus every
test exercising them. `PresencePreview` now shows the resolved Discord Application name (see
ticket 14's follow-up) instead of the button.
