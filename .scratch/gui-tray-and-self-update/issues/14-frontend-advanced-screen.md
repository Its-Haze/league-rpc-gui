# 14: Frontend: Advanced screen

**What to build:** The Advanced section. Discord Application ID: a preset picker (from `GetPresets`) plus a custom-entry field. RPC update interval and stats polling interval, with the min/max from `config` enforced in the UI. A debug-logging toggle bound to `Advanced.DebugMode`.

**Blocked by:** 11, 02

**Status:** ready-for-agent

- [ ] Discord App ID: preset dropdown from `GetPresets` and a validated custom field; empty is rejected inline
- [ ] Update-interval and stats-polling-interval inputs clamped to the `config` bounds, with the bound shown
- [ ] `Advanced.DebugMode` toggle; changing it takes effect without a restart (log level follows it)
- [ ] All changes go through `ApplySettings` and surface `Validate()` errors inline
- [ ] Vitest: bounds enforcement, preset-vs-custom switching
