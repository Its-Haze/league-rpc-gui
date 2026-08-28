package discord

import (
	"github.com/its-haze/league-rpc/internal/config"
	"github.com/its-haze/league-rpc/internal/state"
	"github.com/its-haze/league-rpc/pkg/types"
)

// MapStateToPresence converts application state to Discord RPC data
// This is the main routing function that determines which presence builder to use
func MapStateToPresence(st *state.State, cfg *config.Config) *RPCData {
	if st == nil {
		return &RPCData{}
	}

	// Resolve per-mode display overrides once here; builders read the effective
	// toggles off Display.Default, so hand them a config with those filled in.
	if cfg != nil {
		cfg = withResolvedDisplay(st.GameMode, cfg)
	}

	// Route to appropriate builder based on game flow phase
	switch st.GameFlowPhase {
	case types.GameFlowInProgress:
		// In game - check if TFT or regular game
		if st.GameMode == types.GameModeTFT {
			return BuildTFTInGamePresence(st, cfg)
		}
		// Champion/skin/chroma resolve once per match and ChampionID is
		if st.ChampionID == "" {
			return BuildInChampSelectPresence(st, cfg)
		}
		return BuildInGamePresence(st, cfg)

	case types.GameFlowWatching:
		// Spectating someone else's game.
		return BuildSpectatingPresence(st, cfg)

	case types.GameFlowChampSelect, types.GameFlowGameStart:
		// Champion selection phase
		return BuildInChampSelectPresence(st, cfg)

	case types.GameFlowMatchmaking, types.GameFlowCheckedIntoTournament:
		// In queue
		return BuildInQueuePresence(st, cfg)

	case types.GameFlowReadyCheck:
		// Ready check - keep showing queue state
		return BuildInQueuePresence(st, cfg)

	case types.GameFlowLobby:
		// In lobby - check if custom or matchmaking
		if st.IsCustom || st.IsPractice {
			return BuildInCustomLobbyPresence(st, cfg)
		}
		return BuildInLobbyPresence(st, cfg)

	case types.GameFlowNone, types.GameFlowWaitingForStats,
		types.GameFlowPreEndOfGame, types.GameFlowEndOfGame:
		// Idle in client or post-game
		return BuildInClientPresence(st, cfg)

	default:
		// Unknown phase - default to in client
		return BuildInClientPresence(st, cfg)
	}
}

// resolveDisplay returns the effective ShowRank/ShowStats for mode: its
// Display.Modes override when mode is known, otherwise Display.Default.
func resolveDisplay(cfg *config.Config, mode types.GameMode) (showRank, showStats bool) {
	def := cfg.Display.Default
	if !types.ValidGameMode(mode) {
		return def.ShowRank, def.ShowStats
	}
	r := cfg.Display.Resolve(string(mode))
	return r.ShowRank, r.ShowStats
}

// withResolvedDisplay returns a shallow copy of cfg whose Display.Default holds
// the per-mode-resolved toggles, so presence builders need no override logic.
func withResolvedDisplay(mode types.GameMode, cfg *config.Config) *config.Config {
	showRank, showStats := resolveDisplay(cfg, mode)
	resolved := *cfg
	resolved.Display.Default = config.DisplayDefaults{ShowRank: showRank, ShowStats: showStats}
	return &resolved
}

// ShouldClearPresence returns true if presence should be cleared instead of updated
// This happens when the user has --hide-in-client enabled and is idle
func ShouldClearPresence(st *state.State, cfg *config.Config) bool {
	if !cfg.Presence.ShowInClient && st.GameFlowPhase.IsInClient() {
		return true
	}
	return false
}
