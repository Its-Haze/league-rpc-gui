# 11: Frontend shell: navigation, top strip, theme

**What to build:** The app shell. Left sidebar with sections Home / Display / Behavior / Advanced / Help / About. A persistent top strip showing three connection lights (League process, LCU, Discord) driven by `status:changed`, plus a Pause toggle reachable from every screen. Theme support (system / light / dark) persisted to `Config.Theme`, with `system` following the OS. Run a `frontend-design` pass to set the palette (including the accent), type scale, and spacing tokens; the user reviews the accent.

**Blocked by:** 01, 08

**Status:** done

- [x] Sidebar navigation between the six sections; routing is client-side, no full reloads
- [x] Top strip renders the three connection lights from `GetStatus()` and updates on `status:changed`
- [x] Pause toggle in the top strip calls `App.SetPaused` and reflects `App.IsPaused`
- [x] Theme switch writes `Config.Theme`; `system` tracks OS light/dark live
- [x] `frontend-design` pass done: palette, accent, type scale, spacing as the token file from ticket 01, dark-first, restrained
- [x] In-repo primitives (button, toggle, select, dialog, tabs, field) built on Radix, no component-library dependency
- [x] Vitest covers routing state and the theme resolver

**Notes:** Accent is a Discord-adjacent blue (`#3865d1` light / `#5b8cff` dark), chosen
over a Hextech-gold option and the prior placeholder violet after user review. Display,
Behavior, Advanced, and Help render as labeled placeholders pointing at their own tickets
(12-14, 16); Home and About keep what `App.tsx` already had (the settings echo and the
version/update banner, respectively) until tickets 15 and 16 replace them. The Dialog and
Tabs primitives have no caller yet; they're built for the per-mode override screens and
confirmation flows those tickets need. A follow-up review pass fixed a config-write race in
`App.applyPatch` (two in-flight patches could stomp each other via a stale closure), a
theme picker interactive before the first `GetSettings()` resolved, the error banner being
invisible outside the Home screen, a stale-revert on a delayed `SetPaused` rejection, the
`--text-*` scale never being wired into Tailwind's `@theme inline` block, and the danger
button variant hardcoding white text instead of a paired token.

Post-review user feedback: dropped the sidebar's left accent rule (read as templated) in
favor of a plain filled active row, and swapped the sans-serif stack from Segoe UI to a
self-hosted Inter (`@fontsource-variable/inter`, one variable-weight file, no runtime
Google Fonts dependency), matching Blitz.gg's own font choice at the user's request; sidebar
labels also dropped their uppercase/tracked styling to match. The app icon (`build/appicon.png`
-> `windows/icon.ico`) and the system tray icon (`cmd/league-rpc-gui/icons.go`, `tray-icon.png`)
are now both the League "classic" gold-L crest (`league-rpc-linux/assets/league-classic-borderless.jpg`,
resized). `application.Options.Icon` (`app-icon.png`) is also set at runtime, since Windows
only picks up the linked `.ico` resource from a `task build` output, not a `wails3 dev` run.
