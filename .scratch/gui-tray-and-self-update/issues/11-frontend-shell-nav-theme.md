# 11: Frontend shell: navigation, top strip, theme

**What to build:** The app shell. Left sidebar with sections Home / Display / Behavior / Advanced / Help / About. A persistent top strip showing three connection lights (League process, LCU, Discord) driven by `status:changed`, plus a Pause toggle reachable from every screen. Theme support (system / light / dark) persisted to `Config.Theme`, with `system` following the OS. Run a `frontend-design` pass to set the palette (including the accent), type scale, and spacing tokens; the user reviews the accent.

**Blocked by:** 01, 08

**Status:** ready-for-agent

- [ ] Sidebar navigation between the six sections; routing is client-side, no full reloads
- [ ] Top strip renders the three connection lights from `GetStatus()` and updates on `status:changed`
- [ ] Pause toggle in the top strip calls `App.SetPaused` and reflects `App.IsPaused`
- [ ] Theme switch writes `Config.Theme`; `system` tracks OS light/dark live
- [ ] `frontend-design` pass done: palette, accent, type scale, spacing as the token file from ticket 01, dark-first, restrained
- [ ] In-repo primitives (button, toggle, select, dialog, tabs, field) built on Radix, no component-library dependency
- [ ] Vitest covers routing state and the theme resolver
