# 16: Frontend: Help and About

**What to build:** The Help section: a live log viewer (initial fill from `GetRecentLogs()`, then appended from `log:line` events), an "Open logs folder" button, a "Copy diagnostics" button (app version, Windows build, League / LCU / Discord state, last error) that copies a paste-ready block, and links to the Discord server (`https://discord.haze.sh`), a bug report, and a feature request. Wire the `internal/logging` subscribe path to `log:line` events here. The About section: version from `App.GetVersion()`, "Check for updates", and the rendered changelog from ticket 09. Also needs small `internal/app` bindings for opening the logs folder and building the diagnostics blob.

**Blocked by:** 11, 03, 09

**Status:** ready-for-agent

- [ ] Log viewer fills from `GetRecentLogs()` and live-tails via `log:line`; scroll-lock when the user scrolls up
- [ ] `cmd/league-rpc-gui` emits `log:line` from the `internal/logging` subscribe channel
- [ ] "Open logs folder" binding opens the `%APPDATA%` logs directory in Explorer
- [ ] "Copy diagnostics" binding returns a markdown block; the button copies it to the clipboard
- [ ] Links: Discord (`https://discord.haze.sh`), bug report, and feature request open in the browser
- [ ] About shows the version, a working "Check for updates", and the markdown changelog (with the offline fallback)
- [ ] Vitest: log-tail append and scroll-lock, diagnostics formatting
