# 05: Presence template engine

**What to build:** A new `internal/presence/template` package that renders a `details`/`state` pair from a template string using plain `{token}` substitution (no logic, not Go `text/template`). Per-context token sets for `in-client`, `champ-select`, `in-game`, `spectating`. Unknown tokens are left literal and reported by a preview call. At render time a token with no data resolves to empty, then runs of whitespace and dangling separators (`-`, `|`, `•`) collapse. Ship default templates that reproduce the current hardcoded presence strings exactly. `discord/mapper.go` renders through this engine. Expose `App.RenderTemplatePreview(context, tmpl, sampleData)`.

**Blocked by:** 02

**Status:** ready-for-agent

- [ ] `internal/presence/template` renders `{token}` substitution against a per-context data map
- [ ] Known token set defined per context (`in-client`, `champ-select`, `in-game`, `spectating`)
- [ ] Unknown token: left literal in output, returned in a `[]string` of warnings from the preview API
- [ ] Empty-value token: removed, then whitespace runs collapse to one space and dangling `-`/`|`/`•` are trimmed
- [ ] `config` ships default `Presence.Templates` entries that render to today's exact strings for each context
- [ ] `discord/mapper.go` produces presence text only through this engine
- [ ] `App.RenderTemplatePreview` binding returns rendered `details`/`state` plus warnings
- [ ] Tests: substitution, unknown-token reporting, collapse rules, each context default matches the current string
