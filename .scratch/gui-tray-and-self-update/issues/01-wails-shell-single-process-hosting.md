# 01: Wails shell hosting the daemon in one process

**What to build:** A new entry point `cmd/league-rpc-gui/main.go` that builds the logger and `config.Store`, calls `daemon.Wire(store, logger)`, runs `Daemon.Run(ctx)` in a goroutine, and starts a Wails v3 app whose window loads the frontend and whose bindings are backed by `internal/app`. A hand-assembled `frontend/` (Vite + React + TypeScript + Tailwind v4, Radix primitives available, no component-library dependency) that builds and loads. `cmd/league-rpc` is left untouched as the headless build.

**Blocked by:** None (can start immediately)

**Status:** done

- [x] `cmd/league-rpc-gui/main.go` starts the daemon goroutine and the Wails v3 event loop on the main thread; canceling the context on app exit stops `Daemon.Run` (`OnShutdown` cancels the ctx and waits on `daemonDone`)
- [x] `frontend/` builds with Vite, renders a placeholder React app, and is embedded in the binary (`frontend/embed.go`, `//go:embed all:dist`; `npm run build` and `wails3 build` both produce `dist`)
- [x] Tailwind v4 is wired (`@tailwindcss/vite`, `@import "tailwindcss"`) and `frontend/src/tokens.css` holds the colors/spacing/type CSS variables, theme-aware, provisional values
- [x] `internal/app` gains no Wails import; `cmd/league-rpc-gui/service.go` adapts it, and `publishConfigChanges` bridges `App.SubscribeSettings()` to the `settings:changed` Wails event
- [x] `GetSettings` / `ApplySettings` / `GetPresets` are callable from the frontend (generated bindings) and round-trip a value (`service_test.go`, and `App.tsx` toggles `presence.show_emojis`)
- [~] `wails3 dev` against a real League/Discord: `wails3 build` runs the whole toolchain green (icons, bindings, npm, vite, syso, `go build -tags production`) and emits `bin/league-rpc-gui.exe`; `wails3 dev` with a live client is a manual check on a desktop session
- [x] `go build ./...` and `go test ./...` pass; `cmd/league-rpc` still runs headless (unchanged)
