# League RPC (Go rewrite)

Discord Rich Presence for League of Legends, driven by the LCU (League Client Update) API.

## Language

**State**:
The application's live snapshot of what the summoner is doing right now (client phase, queue, champion/skin, ranked stats).
_Avoid_: ClientData, session

**GameFlowPhase**:
The LCU's own phase string for what the client is doing (`None`, `Lobby`, `ChampSelect`, `InProgress`, `Watching`, ...). Values must mirror the LCU API's raw strings exactly: this is not a value we get to invent, it's a contract with Riot's client.
_Avoid_: GameState, Phase

**Watching**:
The `GameFlowPhase` value meaning the summoner is spectating a game, not playing in one. Presence for this phase shows the spectated match's leading player's champion/skin (or the map icon, for modes without champions) rather than the summoner's own KDA.

**TFT queue variants**:
Ranked TFT (queue ID 1100), Double Up (1160), and Hyper Roll (1130) are each distinct from normal/unranked TFT (1090) and must not share one `QueueID` constant: ranked-stat lookups need the specific ID, not just "is this TFT."

**Updater**:
The single component responsible for every Discord presence send. It debounces rapid `State` changes into one send, then keeps resending periodically (the "heartbeat") for as long as that presence is current, with a faster resend burst (the "reclaim") right after each real change. It exists because Discord lets League's own native Rich Presence integration silently overwrite ours if we stop sending. See [ADR-0001](./docs/adr/0001-updater-owns-the-heartbeat.md).
_Avoid_: RPCUpdater, presence manager

**Live Client Data**:
The unauthenticated local API at `127.0.0.1:2999`, distinct from the LCU: fixed port, no lockfile auth, and only reachable while a match is actually in progress. Used for in-game KDA/champion/gold polling. Owned by its own package (`internal/livegame`), not `internal/lcu`.
_Avoid_: LCU, game data (bare)

**Champion Data**:
Static champion/skin/chroma lookups from Data Dragon and Community Dragon: plain HTTP, no League process required at all. Owned by its own package (`internal/championdata`), not `internal/lcu`.
_Avoid_: LCU, ddragon (as a stand-in for the whole concept)

**League Classic**:
Display name for the `JADE` game mode (raw LCU/game-data identifier), played on the Classic Rift map (map ID 453).

**ARAM: Mayhem**:
Display name for the `KIWI` game mode.

**ARAM: Mayhem Classic-ish**:
Display name for the `KIWI_JADE` game mode: ARAM Mayhem played on the League Classic map. The "-ish" is intentional, not a typo to silently fix; it's the current player-facing string.

**Locale**:
The League client's own display language (e.g. `en_US`), read by inspecting the running League process's `--locale=` command-line argument. This is auto-detected from the OS process, never a user-facing setting in `Config`.
_Avoid_: language, config.Locale

**Daemon**:
The long-running background process itself, launched at Windows startup or manually. It lives in the system tray and keeps running until it crashes or the user explicitly quits it from the tray. Opening/closing the GUI window only shows/hides it; it never starts or stops the daemon. See [ADR-0002](./docs/adr/0002-daemon-never-exits-on-connection-loss.md).
_Avoid_: app, process, service

**Connection Supervisor**:
The component that owns retrying a connection (to Discord, or to the LCU) forever. A closed League client or a closed Discord client is expected: it resumes the retry loop, never exits the Daemon. No LCU connection means no Discord presence, cleared immediately with no grace period. The one exception: while League's process is running but the LCU hasn't connected yet, the placeholder rotation shows instead (see `Updater`). The Discord Connection Supervisor is additionally gated on League's own process: it never connects to Discord, and disconnects if already connected, whenever League isn't running. See [ADR-0003](./docs/adr/0003-discord-connection-gated-by-league.md).
_Avoid_: reconnect logic, watchdog
