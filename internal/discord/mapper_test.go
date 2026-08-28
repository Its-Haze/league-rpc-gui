package discord

import (
	"testing"

	"github.com/its-haze/league-rpc/internal/config"
	"github.com/its-haze/league-rpc/internal/state"
	"github.com/its-haze/league-rpc/pkg/types"
)

func TestResolveDisplay(t *testing.T) {
	tr, fa := true, false

	tests := []struct {
		name              string
		def               config.DisplayDefaults
		modes             map[string]config.ModeOverride
		mode              types.GameMode
		wantRank, wantSts bool
	}{
		{
			name:     "no override inherits the default",
			def:      config.DisplayDefaults{ShowRank: true, ShowStats: true},
			mode:     types.GameModeClassic,
			wantRank: true, wantSts: true,
		},
		{
			name:     "override wins over the default",
			def:      config.DisplayDefaults{ShowRank: true, ShowStats: true},
			modes:    map[string]config.ModeOverride{string(types.GameModeArena): {ShowRank: &fa}},
			mode:     types.GameModeArena,
			wantRank: false, wantSts: true,
		},
		{
			name:     "override can turn a toggle on",
			def:      config.DisplayDefaults{ShowRank: false, ShowStats: false},
			modes:    map[string]config.ModeOverride{string(types.GameModeARAM): {ShowStats: &tr}},
			mode:     types.GameModeARAM,
			wantRank: false, wantSts: true,
		},
		{
			name:     "unknown mode falls back to the default and ignores any stored override",
			def:      config.DisplayDefaults{ShowRank: true, ShowStats: false},
			modes:    map[string]config.ModeOverride{"NOT_A_MODE": {ShowRank: &fa}},
			mode:     "NOT_A_MODE",
			wantRank: true, wantSts: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.DefaultConfig()
			cfg.Display.Default = tt.def
			cfg.Display.Modes = tt.modes

			gotRank, gotSts := resolveDisplay(cfg, tt.mode)
			if gotRank != tt.wantRank || gotSts != tt.wantSts {
				t.Errorf("resolveDisplay() = (%v, %v), want (%v, %v)", gotRank, gotSts, tt.wantRank, tt.wantSts)
			}
		})
	}
}

func TestMapStateToPresence_HonorsPerModeRankOverride(t *testing.T) {
	fa := false

	inGame := func() *state.State {
		st := state.NewState()
		st.GameFlowPhase = types.GameFlowInProgress
		st.GameMode = types.GameModeClassic
		st.ChampionID = "Ahri"
		st.ChampionName = "Ahri"
		st.QueueID = types.QueueSoloQ
		st.SummonerRank = types.RankedStats{Tier: types.TierGold, Division: types.DivisionIV, LeaguePoints: 40}
		return st
	}

	// Default hides rank; a CLASSIC override turns it back on.
	cfg := config.DefaultConfig()
	cfg.Display.Default.ShowRank = false
	cfg.Display.Modes = map[string]config.ModeOverride{
		string(types.GameModeClassic): {ShowRank: boolPtr(true)},
	}
	if got := MapStateToPresence(inGame(), cfg); got.SmallText != "Gold IV: 40 LP" {
		t.Errorf("override ShowRank=true not reflected: SmallText = %q", got.SmallText)
	}

	// Default shows rank; a CLASSIC override hides it.
	cfg = config.DefaultConfig()
	cfg.Display.Default.ShowRank = true
	cfg.Display.Modes = map[string]config.ModeOverride{
		string(types.GameModeClassic): {ShowRank: &fa},
	}
	if got := MapStateToPresence(inGame(), cfg); got.SmallText == "Gold IV: 40 LP" {
		t.Error("override ShowRank=false not reflected: rank still shown")
	}
}

func boolPtr(b bool) *bool { return &b }

func TestPerModeDisplayApplies(t *testing.T) {
	tests := []struct {
		phase types.GameFlowPhase
		want  bool
	}{
		{types.GameFlowNone, false},
		{types.GameFlowWaitingForStats, false},
		{types.GameFlowPreEndOfGame, false},
		{types.GameFlowEndOfGame, false},
		{types.GameFlowLobby, true},
		{types.GameFlowMatchmaking, true},
		{types.GameFlowChampSelect, true},
		{types.GameFlowInProgress, true},
		{types.GameFlowWatching, true},
	}
	for _, tt := range tests {
		if got := perModeDisplayApplies(tt.phase); got != tt.want {
			t.Errorf("perModeDisplayApplies(%q) = %v, want %v", tt.phase, got, tt.want)
		}
	}
}

// A per-mode override for the last-played mode must not bleed into the
// post-game card, which keeps GameMode from the finished match.
func TestMapStateToPresence_PostGameIgnoresStaleModeOverride(t *testing.T) {
	postGame := func() *state.State {
		st := state.NewState()
		st.GameFlowPhase = types.GameFlowEndOfGame
		st.GameMode = types.GameModeARAM // left over from the match that just ended
		return st
	}

	base := config.DefaultConfig()
	withOverride := config.DefaultConfig()
	withOverride.Display.Modes = map[string]config.ModeOverride{
		string(types.GameModeARAM): {ShowRank: boolPtr(false), ShowStats: boolPtr(false)},
	}

	got := MapStateToPresence(postGame(), withOverride)
	want := MapStateToPresence(postGame(), base)
	if !got.Equals(want) {
		t.Errorf("post-game presence changed by an ARAM override:\n got %+v\nwant %+v", *got, *want)
	}
}
