# 12: Frontend: Display screen

**What to build:** The Display section. Global toggles for show rank, show stats (KDA/CS), and status emojis. A collapsible list of `GameMode`s, each with a "use default / override" control for rank and stats. Per-context template editors (in client, champ select, in game, spectating) for the `details` and `state` lines, each with a live preview rendered through `App.RenderTemplatePreview` against sample data, showing unknown-token warnings. All changes go through `ApplySettings`.

**Blocked by:** 11, 05, 06

**Status:** done

- [x] Global show-rank / show-stats / show-emojis toggles bound to the config tree
- [x] Per-`GameMode` rows with default-or-override for rank and stats; mode list comes from a binding, not hardcoded
- [x] Template editors for all four contexts, `details` and `state` each
- [x] Live preview per editor via `App.RenderTemplatePreview`, updating as the user types, listing unknown-token warnings
- [x] A token reference (available `{tokens}` for the current context) is visible near each editor
- [x] Invalid config (e.g. empty required field) is shown inline and blocks save, matching `Validate()`
- [x] Vitest: override-resolution helper, editor state, preview wiring with a mocked binding

**Notes:** Added two small `internal/app` bindings this ticket needed that didn't exist yet:
`GetGameModes()` (wraps `pkg/types.GameModes()`, so the mode list can't drift from the const
block the way ticket 19 warned about) and `GetTemplateTokens(ctx)` (wraps
`template.KnownTokens`) for the per-editor token reference. Per-mode overrides render as a
default/on/off `Select` per field (`ModeOverrideRow`); "default" clears the field back to nil
rather than writing an explicit value. The four context editors (`TemplateEditor`) call
`RenderTemplatePreview` on every keystroke against `sample=null`, so the preview always matches
what the daemon would render for a blank line. The `Validate()`-inline case that applies here is
an empty `discord_app_id`, which lives on the Advanced screen, not Display; Display's own fields
have no required-non-empty constraint in `Validate()`, so the inline-error path here is just the
shared `useSettings` error banner from a rejected `ApplySettings`.
