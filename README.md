# League RPC (Go Rewrite)

A modern, fast, and user-friendly Discord Rich Presence application for League of Legends, built with Go and Wails.

## Overview

League RPC is designed to provide:
- **Better Performance**: Native Go application with minimal resource usage
- **Modern GUI**: Clean, intuitive interface built with React + shadcn/ui
- **Easy Configuration**: All settings accessible through the GUI (no command-line flags!)
- **Rich Features**: Full Discord Rich Presence coverage for every phase of a League session

## Features

- **Real-time Discord Rich Presence** with summoner icons and rank emblems
- **In-game Stats Display**: KDA, Creep Score, and game timer
- **Champion Skins**: Display selected skins (including animated Ultimate skins)
- **Rank Display**: Show your rank across all game modes (SoloQ, Flex, TFT, Arena)
- **TFT Support**: Display your favorite TFT companion
- **Easy Configuration**: All settings managed through the GUI
- **Auto-update Checker**: Get notified of new versions

## Technology stack

### Backend (Go)

See `go.mod` for the pinned versions.

| Dependency | Purpose |
|------------|---------|
| **lcu-gopher** | League Client API integration (WebSocket + HTTP) |
| **internal/discord/ipc** | Discord Rich Presence transport (in-house, see DEPENDENCIES.md) |
| **gopsutil/v4** | Process detection (Discord, League) |
| **zerolog** | High-performance structured logging |

### Frontend (To be added)

- **Wails v3**: Desktop application framework
- **React + TypeScript**: Modern UI framework
- **shadcn/ui**: Beautiful, accessible components
- **Tailwind CSS**: Utility-first styling

## Development status

Currently in early development.

### Completed
- [x] Project initialization
- [x] Core dependencies installed
- [x] Directory structure created
- [x] Architecture designed

### In Progress
- [ ] Configuration system
- [ ] LCU client wrapper
- [ ] Discord RPC integration
- [ ] State management
- [ ] Frontend setup with Wails

### Planned
- [ ] Settings GUI
- [ ] System tray integration
- [ ] Auto-launch functionality
- [ ] Update checker
- [ ] Comprehensive testing

## Configuration

All settings are stored in `%APPDATA%\league-rpc\config.json` and managed through the GUI:

- **Discord App ID**: Choose from presets or use custom
- **Display Options**: Toggle stats, rank, emojis, in-client presence
- **League Settings**: Auto-launch League, custom installation path
- **Advanced**: Update interval, debug mode

## Building from Source

### Prerequisites
- Go 1.22 or later
- Node.js 18+ (for frontend)
- Wails CLI: `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`

### Development

```bash
# Clone the repository
git clone https://github.com/its-haze/league-rpc
cd league-rpc

# Install dependencies
go mod download

# Run in development mode (after Wails setup)
wails3 dev

# Build for production
wails3 build
```

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
