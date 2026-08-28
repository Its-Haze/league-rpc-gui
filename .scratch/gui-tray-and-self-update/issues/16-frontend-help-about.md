# 16: Frontend: Help and About

**What to build:** The Help section: a live log viewer (initial fill from `GetRecentLogs()`, then appended from `log:line` events), an "Open logs folder" button, a "Copy diagnostics" button (app version, Windows build, League / LCU / Discord state, last error) that copies a paste-ready block, and links to the Discord server (`https://discord.haze.sh`), a bug report, and a feature request. Wire the `internal/logging` subscribe path to `log:line` events here. The About section: version from `App.GetVersion()`, "Check for updates", and the rendered changelog from ticket 09. Also needs small `internal/app` bindings for opening the logs folder and building the diagnostics blob.

**Blocked by:** 11, 03, 09

**Status:** done

- [x] Log viewer fills from `GetRecentLogs()` and live-tails via `log:line`; scroll-lock when the user scrolls up
- [x] `cmd/league-rpc-gui` emits `log:line` from the `internal/logging` subscribe channel
- [x] "Open logs folder" binding opens the `%APPDATA%` logs directory in Explorer
- [x] "Copy diagnostics" binding returns a markdown block; the button copies it to the clipboard
- [x] Links: Discord (`https://discord.haze.sh`), bug report, and feature request open in the browser
- [x] About shows the version, a working "Check for updates", and the markdown changelog (with the offline fallback)
- [x] Vitest: log-tail append and scroll-lock, diagnostics formatting

**Notes:** Added the `internal/app` bindings this ticket needed: `GetRecentLogs`/`SubscribeLogs`
(wrap the existing `logging.Ring`), `OpenLogsFolder` (shells to Explorer via an injectable
`openFolder func(string) error`, so tests never actually shell out), and `GetDiagnostics`
(version, OS, the three connection states, paused, phase, last update error). `SubscribeLogs`
is bridged to a `log:line` Wails event by a new `guiService.publishLogLines` goroutine, wired in
`main.go` next to the existing `publishConfigChanges`. "Copy diagnostics" wraps the backend
string in a Markdown code fence (`lib/diagnostics.ts`) before writing it to the clipboard, so
pasting into a GitHub issue renders as a preformatted block. Bug-report and feature-request URLs
(`lib/links.ts`) point at `github.com/its-haze/league-rpc/issues/new?template=...`, matching the
module path; ticket 17's issue templates and the eventual canonical repo URL are unverified
against this yet. About keeps `UpdateBanner` for the actionable "a new version is available" banner
and adds its own persistent version/check-for-updates/changelog block below it, since `UpdateBanner`
hides itself entirely when no update is available and the ticket wants those controls visible
regardless.
