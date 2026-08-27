# 12: Frontend: Display screen

**What to build:** The Display section. Global toggles for show rank, show stats (KDA/CS), and status emojis. A collapsible list of `GameMode`s, each with a "use default / override" control for rank and stats. Per-context template editors (in client, champ select, in game, spectating) for the `details` and `state` lines, each with a live preview rendered through `App.RenderTemplatePreview` against sample data, showing unknown-token warnings. All changes go through `ApplySettings`.

**Blocked by:** 11, 05, 06

**Status:** ready-for-agent

- [ ] Global show-rank / show-stats / show-emojis toggles bound to the config tree
- [ ] Per-`GameMode` rows with default-or-override for rank and stats; mode list comes from a binding, not hardcoded
- [ ] Template editors for all four contexts, `details` and `state` each
- [ ] Live preview per editor via `App.RenderTemplatePreview`, updating as the user types, listing unknown-token warnings
- [ ] A token reference (available `{tokens}` for the current context) is visible near each editor
- [ ] Invalid config (e.g. empty required field) is shown inline and blocks save, matching `Validate()`
- [ ] Vitest: override-resolution helper, editor state, preview wiring with a mocked binding
