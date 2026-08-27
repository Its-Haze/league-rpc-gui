package discord

import (
	"testing"
	"time"

	"github.com/its-haze/league-rpc/internal/config"
	"github.com/its-haze/league-rpc/internal/state"
	"github.com/its-haze/league-rpc/pkg/constants"
	"github.com/its-haze/league-rpc/pkg/types"
)

func TestFormatSkinName(t *testing.T) {
	tests := []struct {
		name                                     string
		championName, skinName, chromaName, want string
	}{
		{"skin found: skin name alone, no champion prefix", "Qiyana", "Battle Boss Qiyana", "", "Battle Boss Qiyana"},
		{"chroma found: skin name plus chroma, no champion prefix", "Chogath", "Battlecast Cho'Gath", "Emerald", "Battlecast Cho'Gath (Emerald)"},
		{"default skin: falls back to champion name", "Ahri", "", "", "Ahri"},
		{"skin literally \"default\": falls back to champion name", "Ahri", "default", "", "Ahri"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatSkinName(tt.championName, tt.skinName, tt.chromaName); got != tt.want {
				t.Errorf("FormatSkinName(%q, %q, %q) = %q, want %q", tt.championName, tt.skinName, tt.chromaName, got, tt.want)
			}
		})
	}
}

func TestFormatGameModeName(t *testing.T) {
	tests := []struct {
		mode types.GameMode
		want string
	}{
		{types.GameModeClassic, "Summoner's Rift"},
		{types.GameModeARAM, "Howling Abyss (ARAM)"},
		{types.GameModeBrawl, "Brawl"},
		{types.GameModeArena, "Arena"},
		{types.GameModeSwarm, "Swarm"},
		{types.GameModeUltimateSpellbook, "Ultimate Spellbook"},
		{"SOME_UNKNOWN_MODE", "SOME_UNKNOWN_MODE"},
	}
	for _, tt := range tests {
		if got := FormatGameModeName(tt.mode); got != tt.want {
			t.Errorf("FormatGameModeName(%q) = %q, want %q", tt.mode, got, tt.want)
		}
	}
}

func TestQueueDisplayName(t *testing.T) {
	tests := []struct {
		name string
		st   *state.State
		want string
	}{
		{
			name: "detailed description wins",
			st: &state.State{
				QueueDetailedDescription: "Ranked Solo/Duo",
				QueueName:                "5v5 Ranked Solo games",
			},
			want: "Ranked Solo/Duo",
		},
		{
			name: "falls back to queue name",
			st:   &state.State{QueueName: "Custom Game"},
			want: "Custom Game",
		},
		{
			name: "falls back to game mode display name",
			st:   &state.State{GameMode: types.GameModeARAM},
			want: "Howling Abyss (ARAM)",
		},
		{
			name: "final catch-all",
			st:   &state.State{},
			want: "League of Legends",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := queueDisplayName(tt.st); got != tt.want {
				t.Errorf("queueDisplayName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildInGamePresence_StatePrefixedWithInGame(t *testing.T) {
	cfg := config.DefaultConfig()
	st := state.NewState()
	st.Kills, st.Deaths, st.Assists, st.CreepScore = 3, 2, 5, 120

	cfg.ShowStats = true
	got := BuildInGamePresence(st, cfg)
	if want := "In Game · 3/2/5 · 120cs"; got.State != want {
		t.Errorf("State (stats shown) = %q, want %q", got.State, want)
	}

	cfg.ShowStats = false
	got = BuildInGamePresence(st, cfg)
	if want := "In Game"; got.State != want {
		t.Errorf("State (stats hidden) = %q, want %q", got.State, want)
	}
}

func TestBuildTFTInGamePresence_UsesResolvedGameStartTime(t *testing.T) {
	cfg := config.DefaultConfig()
	st := state.NewState()
	st.GameStartTime = 1000

	got := BuildTFTInGamePresence(st, cfg)
	if got.Start != 1000 {
		t.Errorf("Start = %d, want 1000 (the poller-supplied GameStartTime)", got.Start)
	}
}

func TestBuildTFTInGamePresence_FallsBackToNowWhenGameStartTimeUnresolved(t *testing.T) {
	cfg := config.DefaultConfig()
	st := state.NewState()
	st.GameStartTime = 0

	before := time.Now().Unix()
	got := BuildTFTInGamePresence(st, cfg)
	if got.Start < before {
		t.Errorf("Start = %d, want >= %d (fallback to time.Now())", got.Start, before)
	}
}

func TestBuildInCustomLobbyPresence_DetailsIsRawQueueNameStateIsFixed(t *testing.T) {
	cfg := config.DefaultConfig()
	st := state.NewState()
	st.QueueName = "Practice Tool"
	st.IsPractice = true

	got := BuildInCustomLobbyPresence(st, cfg)
	if got.Details != "Practice Tool" {
		t.Errorf("Details = %q, want %q (no prefix)", got.Details, "Practice Tool")
	}
	// State is always the fixed "In Lobby" string, never a repeat of
	// Details/queue name, regardless of custom vs. practice tool.
	if got.State != "In Lobby" {
		t.Errorf("State = %q, want %q", got.State, "In Lobby")
	}
}

func TestBuildInLobbyPresence_RankKnownSwapsLargeTextToCredit(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.ShowRank = true
	st := state.NewState()
	st.QueueID = types.QueueSoloQ
	st.SummonerRank = types.RankedStats{Tier: types.TierGold, Division: types.DivisionIV, LeaguePoints: 40}

	got := BuildInLobbyPresence(st, cfg)
	if got.LargeText != constants.SmallText {
		t.Errorf("LargeText = %q, want the credit line %q", got.LargeText, constants.SmallText)
	}
	if got.SmallText != "Gold IV: 40 LP" {
		t.Errorf("SmallText = %q, want %q", got.SmallText, "Gold IV: 40 LP")
	}
}

func TestMapStateToPresence_InProgressWithUnresolvedChampionFallsBackToChampSelect(t *testing.T) {
	cfg := config.DefaultConfig()
	st := state.NewState()
	st.GameFlowPhase = types.GameFlowInProgress
	st.GameMode = types.GameModeClassic
	// ChampionID left empty: this match's poller hasn't resolved a champion
	// yet (e.g. right after entering InProgress clears the previous match's data).

	got := MapStateToPresence(st, cfg)
	want := BuildInChampSelectPresence(st, cfg)
	if got.State != want.State || got.Details != want.Details {
		t.Errorf("MapStateToPresence() = %+v, want champ-select presence %+v", got, want)
	}
}

func TestMapStateToPresence_InProgressWithResolvedChampionUsesInGamePresence(t *testing.T) {
	cfg := config.DefaultConfig()
	st := state.NewState()
	st.GameFlowPhase = types.GameFlowInProgress
	st.GameMode = types.GameModeClassic
	st.ChampionID = "Chogath"
	st.ChampionName = "Cho'Gath"

	got := MapStateToPresence(st, cfg)
	if got.LargeText != "Cho'Gath" {
		t.Errorf("LargeText = %q, want the resolved champion's display name", got.LargeText)
	}
}

func TestBuildInGamePresence_RankKnownDoesNotSwapLargeText(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.ShowRank = true
	st := state.NewState()
	st.ChampionName = "Ahri"
	st.QueueID = types.QueueSoloQ
	st.SummonerRank = types.RankedStats{Tier: types.TierGold, Division: types.DivisionIV, LeaguePoints: 40}

	got := BuildInGamePresence(st, cfg)
	if got.LargeText == constants.SmallText {
		t.Errorf("LargeText should stay the skin name, not the credit line, got %q", got.LargeText)
	}
}
