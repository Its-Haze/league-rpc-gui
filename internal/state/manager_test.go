package state

import (
	"testing"

	"github.com/its-haze/league-rpc/pkg/types"
	"github.com/rs/zerolog"
)

func TestManager_ClearInGameData_ResetsMatchDataButNotOtherState(t *testing.T) {
	m := NewManager(zerolog.Nop())
	m.UpdateChampion("Chogath", "Cho'Gath", "Battlecast Cho'Gath", "", 1)
	m.UpdateGameStart(1000)
	m.UpdateInGameStats(3, 2, 5, 120, 9, 1500)
	m.UpdateQueue("Ranked Solo/Duo", "RANKED_SOLO_5x5", "", 420, true)

	m.ClearInGameData()

	got := m.Get()
	if got.ChampionID != "" || got.ChampionName != "" || got.SkinName != "" || got.SkinID != 0 || got.GameStartTime != 0 {
		t.Errorf("ClearInGameData() left champion fields set: %+v", got)
	}
	if got.Kills != 0 || got.Deaths != 0 || got.Assists != 0 || got.CreepScore != 0 || got.Level != 0 || got.Gold != 0 {
		t.Errorf("ClearInGameData() left KDA/CS/level/gold set: %+v", got)
	}
	if got.QueueName != "Ranked Solo/Duo" {
		t.Errorf("ClearInGameData() must not touch unrelated fields, got QueueName = %q", got.QueueName)
	}
}

func TestManager_UpdateGameFlowPhase_EnteringInProgressClearsPreviousMatchData(t *testing.T) {
	m := NewManager(zerolog.Nop())
	m.UpdateGameFlowPhase(string(types.GameFlowInProgress))
	m.UpdateChampion("Ahri", "Ahri", "", "", 0)
	m.UpdateInGameStats(0, 1, 0, 120, 0, 0)

	// Match ends, then a brand new match starts: entering InProgress again
	m.UpdateGameFlowPhase(string(types.GameFlowEndOfGame))
	m.UpdateGameFlowPhase(string(types.GameFlowChampSelect))
	m.UpdateGameFlowPhase(string(types.GameFlowInProgress))

	got := m.Get()
	if got.ChampionID != "" {
		t.Errorf("ChampionID = %q, want cleared on entering a new match", got.ChampionID)
	}
	if got.Kills != 0 || got.Deaths != 0 || got.CreepScore != 0 {
		t.Errorf("KDA/CS not cleared on entering a new match: %+v", got)
	}
}

func TestManager_UpdateGameFlowPhase_StayingInProgressDoesNotClearData(t *testing.T) {
	// A reconnect blip re-fetches the current phase unconditionally; if it's
	// still InProgress, an already-resolved mid-match champion must survive.
	m := NewManager(zerolog.Nop())
	m.UpdateGameFlowPhase(string(types.GameFlowInProgress))
	m.UpdateChampion("Ahri", "Ahri", "", "", 0)

	m.UpdateGameFlowPhase(string(types.GameFlowInProgress))

	if got := m.Get().ChampionID; got != "Ahri" {
		t.Errorf("ChampionID = %q, want Ahri preserved across a same-phase update", got)
	}
}

func TestState_Equals_DetectsTFTCompanionDescriptionChange(t *testing.T) {
	a := NewState()
	a.TFTCompanionDescription = "one"
	b := a.Copy()
	b.TFTCompanionDescription = "two"

	if a.Equals(b) {
		t.Error("Equals() = true for states differing only in TFTCompanionDescription")
	}
}
