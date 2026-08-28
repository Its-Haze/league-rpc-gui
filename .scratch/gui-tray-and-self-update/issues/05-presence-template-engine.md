# 05: Presence template engine

**What to build:** A new `internal/presence/template` package that renders a `details`/`state` pair from a template string using plain `{token}` substitution (no logic, not Go `text/template`). Per-context token sets for `in-client`, `champ-select`, `in-game`, `spectating`. Unknown tokens are left literal and reported by a preview call. At render time a token with no data resolves to empty, then runs of whitespace and dangling separators (`-`, `|`, `•`) collapse. Ship default templates that reproduce the current hardcoded presence strings exactly. `discord/mapper.go` renders through this engine. Expose `App.RenderTemplatePreview(context, tmpl, sampleData)`.

**Blocked by:** 02

**Status:** done

- [x] `internal/presence/template` renders `{token}` substitution against a per-context data map
- [x] Known token set defined per context (`in-client`, `champ-select`, `in-game`, `spectating`)
- [x] Unknown token: left literal in output, returned in a `[]string` of warnings from the preview API
- [x] Empty-value token: removed, then whitespace runs collapse to one space and dangling `-`/`|`/`•` are trimmed
- [x] `config` ships default `Presence.Templates` entries that render to today's exact strings for each context
- [x] `discord/mapper.go` produces presence text only through this engine
- [x] `App.RenderTemplatePreview` binding returns rendered `details`/`state` plus warnings
- [x] Tests: substitution, unknown-token reporting, collapse rules, each context default matches the current string

## Comments

### Implementation notes

- **Scope of "renders through this engine":** the four contexts the spec names (`in-client`, `champ-select`, `in-game`, `spectating`) now build their `Details`/`State` through `renderPresenceText`. The lobby, queue, custom-lobby and TFT builders are untouched, matching the spec's out-of-scope line ("any presence-mapping behavior change beyond routing text through the template engine"). Images, `Start`, and `LargeText` stay computed in the builders.
- **Separator set includes U+00B7.** The spec lists `-`, `|`, `•` (U+2022). This app's own strings use the middot U+00B7 ("In Game · {stats}", `FormatKDA`, the arena/swarm lines), so the engine trims it too. Without it the default in-game state renders `In Game ·` when stats are hidden instead of `In Game`. `separators = "-|•·"`.
- **Collapse is seam-local, not global.** Only the whitespace around an empty token collapses; literal whitespace an author typed is preserved. This is what lets the in-client default keep its two spaces (`{emoji}  {availability}` -> `🟢  Online`) while still cleaning up to `Online` when the emoji token is empty.
- **Empty user template falls back to the default.** `renderPresenceText` treats an empty `Details`/`State` override as "unset" and uses the built-in. A user can't blank a line to nothing (Discord rejects empty details anyway).
- **`clamp()` backfills missing template keys** so a migrated v1 config and the GUI always have all four entries to show; it never overwrites a key the user set.
- `App.RenderTemplatePreview(ctx string, tmpl config.TemplatePair, sample map[string]string)` returns `app.TemplatePreview{Details, State, Warnings}`. Empty (nil or `{}`) `sample` uses `template.SampleData(ctx)`. Unknown context is an error; unknown tokens are warnings formatted `unknown token {name}`. Passed through `guiService` and into the regenerated frontend bindings.
- **Preview and runtime share one path.** Both `discord.renderPresenceText` and `RenderTemplatePreview` call `template.RenderPair(ctx, overrideDetails, overrideState, data)`, which applies the blank-line -> default fallback before rendering. A blank field previews as the default it will actually send.

### Code review follow-up (medium)

- **Fixed:** preview skipped the blank-line -> default fallback (preview could differ from runtime); a non-nil empty `sample` map bypassed the sample-data guard and rendered every token empty; four doc comments were truncated by the comment-cap hook.
- **Deferred to ticket 12:** a typo'd token in a saved template is not rejected by `Config.Validate` and renders literally to Discord. This matches the spec ("unknown token: left literal, reported by the preview API"), so there is no save-time rejection by design. `template.KnownTokens(ctx)` is the hook the settings screen uses to warn as the user types.
