package lcu

import (
	"context"
	"errors"
	"testing"

	"github.com/its-haze/league-rpc/internal/state"
	"github.com/its-haze/league-rpc/pkg/types"
	"github.com/rs/zerolog"
)

// fakeGameModeFetcher lets tests control what Live Client Data reports
// without a real League client.
type fakeGameModeFetcher struct {
	mode string
	err  error
}

func (f *fakeGameModeFetcher) GameMode(ctx context.Context) (string, error) {
	return f.mode, f.err
}

func newTestClient(liveGame GameModeFetcher) (*Client, *state.Manager) {
	stateMgr := state.NewManager(zerolog.Nop())
	c := &Client{state: stateMgr, logger: zerolog.Nop(), liveGame: liveGame}
	return c, stateMgr
}

func TestRecoverGameModeFromLiveClientData_PracticeTool(t *testing.T) {
	c, stateMgr := newTestClient(&fakeGameModeFetcher{mode: "PRACTICETOOL"})
	stateMgr.UpdateGameFlowPhase(string(types.GameFlowInProgress))

	c.recoverGameModeFromLiveClientData()

	got := stateMgr.Get()
	if got.GameMode != "PRACTICETOOL" {
		t.Errorf("GameMode = %q, want PRACTICETOOL", got.GameMode)
	}
	if got.QueueName != "Practice Tool" {
		t.Errorf("QueueName = %q, want Practice Tool", got.QueueName)
	}
	if !got.IsPractice {
		t.Error("IsPractice = false, want true")
	}
}

func TestRecoverGameModeFromLiveClientData_NonPracticeDoesNotSetQueueName(t *testing.T) {
	c, stateMgr := newTestClient(&fakeGameModeFetcher{mode: "CLASSIC"})
	stateMgr.UpdateGameFlowPhase(string(types.GameFlowInProgress))

	c.recoverGameModeFromLiveClientData()

	got := stateMgr.Get()
	if got.GameMode != "CLASSIC" {
		t.Errorf("GameMode = %q, want CLASSIC", got.GameMode)
	}
	if got.QueueName != "" || got.IsPractice {
		t.Errorf("expected no queue name/practice flag for a non-practice mode, got %+v", got)
	}
}

func TestRecoverGameModeFromLiveClientData_NoopWhenNotInProgress(t *testing.T) {
	c, stateMgr := newTestClient(&fakeGameModeFetcher{mode: "PRACTICETOOL"})
	stateMgr.UpdateGameFlowPhase(string(types.GameFlowChampSelect))

	c.recoverGameModeFromLiveClientData()

	if got := stateMgr.Get().GameMode; got != types.GameModeClassic {
		t.Errorf("GameMode = %q, want unchanged default (not in progress yet)", got)
	}
}

func TestRecoverGameModeFromLiveClientData_NoopOnFetchError(t *testing.T) {
	c, stateMgr := newTestClient(&fakeGameModeFetcher{err: errors.New("live client unreachable")})
	stateMgr.UpdateGameFlowPhase(string(types.GameFlowInProgress))

	c.recoverGameModeFromLiveClientData()

	if got := stateMgr.Get().GameMode; got != types.GameModeClassic {
		t.Errorf("GameMode = %q, want unchanged default (fetch failed)", got)
	}
}

func TestRecoverGameModeFromLiveClientData_NoopWhenFetcherNil(t *testing.T) {
	c, stateMgr := newTestClient(nil)
	stateMgr.UpdateGameFlowPhase(string(types.GameFlowInProgress))

	c.recoverGameModeFromLiveClientData() // must not panic on a nil liveGame
}
