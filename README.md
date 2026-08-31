# League RPC (Go Rewrite)

A modern, fast, and user-friendly Discord Rich Presence application for League of Legends, built with Go and Wails.

## Overview

League RPC is designed to provide:
- **Better Performance**: Native Go application with minimal resource usage
- **Modern GUI**: Clean interface built with React, Tailwind CSS, and Radix UI
- **Easy Configuration**: All settings accessible through the GUI (no command-line flags!)
- **Rich Features**: Full Discord Rich Presence coverage for every phase of a League session

## Features

- **Real-time Discord Rich Presence** with summoner icons and rank emblems
- **In-game Stats Display**: KDA, Creep Score, and game timer
- **Champion Skins**: Display selected skins (including animated Ultimate skins)
- **Rank Display**: Show your rank across all game modes (SoloQ, Flex, TFT, Arena)
- **TFT Support**: Display your favorite TFT companion
- **Customizable Presence Text**: edit the details/state shown for each phase (in client, champ select, in game, spectating), with a live preview
- **System Tray**: closing the window hides it to the tray; presence keeps running until you quit from there
- **Start with Windows**: optional, with the window starting hidden
- **In-app Updates**: checks GitHub Releases and installs a signed update with one click, no manual download

## Technology stack

### Backend (Go)

See `go.mod` for the pinned versions.

| Dependency | Purpose |
|------------|---------|
| **lcu-gopher** | League Client API integration (WebSocket + HTTP) |
| **internal/discord/ipc** | Discord Rich Presence transport (in-house, see DEPENDENCIES.md) |
| **gopsutil/v4** | Process detection (Discord, League) |
| **zerolog** + **lumberjack** | Structured logging to a rotating file and an in-memory ring buffer |
| **Wails v3** | Hosts the GUI and the daemon in one process; system tray, window, and the signed in-app updater |

### Frontend

See `frontend/package.json` for pinned versions.

| Dependency | Purpose |
|------------|---------|
| **React + TypeScript** | UI framework |
| **Tailwind CSS + Radix UI primitives** | Styling and accessible component behavior (no shadcn/ui, see DEPENDENCIES.md) |
| **lucide-react** | Icons |
| **dompurify** | Sanitizes the rendered changelog markdown |
| **Vitest + React Testing Library** | Frontend unit tests |

## Development status

The GUI, system tray, config, presence templates, and in-app updates described above are built. See `.scratch/gui-tray-and-self-update/` for the milestone's tickets and history.

### Remaining
- [ ] Windows code-signing certificate (unsigned builds currently trigger a SmartScreen warning; see ADR-0005)
- [ ] macOS / Linux support (not planned; see Platform Support below)

## Configuration

All settings are stored in `%APPDATA%\league-rpc\config.json` and managed through the GUI:

- **Display**: Toggle rank, stats, emojis, and in-client presence; edit the presence text template for each phase (in client, champ select, in game, spectating), with a live preview
- **Behavior**: Start with Windows, what closing the window does, pause presence
- **Advanced**: Discord App ID (presets or custom), update/polling intervals, debug logging

## Building from Source

### Prerequisites
- Go (see `go.mod` for the version)
- Node.js 18+ (for the frontend)
- [Task](https://taskfile.dev/) and the Wails v3 CLI: `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`

### Development

```bash
# Clone the repository
git clone https://github.com/its-haze/league-rpc
cd league-rpc

# Install Go dependencies
go mod download

# Run in development mode, with hot reload
task dev

# Build a binary (bin/league-rpc-gui.exe)
task build

# Build a distributable Windows installer
task package
```

The headless daemon (`cmd/league-rpc`, no GUI, useful for debugging) builds separately: `go build ./cmd/league-rpc`.

## Platform Support

**Windows only** - This rewrite targets Windows exclusively due to Riot's Vanguard anti-cheat breaking Linux support.

## Why Go

1. **Performance**: Go's compiled nature provides instant startup and minimal resource usage
2. **User Experience**: Modern GUI makes configuration effortless
3. **Maintainability**: Clean architecture and type safety make the codebase easier to maintain
4. **Distribution**: Single executable with no external runtime dependencies required

## Credits

- LCU API library: [lcu-gopher](https://github.com/its-haze/lcu-gopher)

## License

MIT License - See [LICENSE](LICENSE) for details
