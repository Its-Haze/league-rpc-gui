package livegame

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/its-haze/league-rpc/internal/championdata"
	"github.com/its-haze/league-rpc/internal/config"
	"github.com/its-haze/league-rpc/pkg/types"
	"github.com/rs/zerolog"
)

const (
	testInterval = 2 * time.Millisecond
	testTimeout  = 2 * time.Second
)

// storeWithInterval builds a config.Store whose StatsPollingInterval equals d.
func storeWithInterval(d time.Duration) *config.Store {
	return config.NewStore(&config.Config{StatsPollingInterval: int(d / time.Millisecond)})
}

// fakeResolver lets tests control what championdata.Resolve returns without
// touching real Data Dragon/Meraki.
type fakeResolver struct {
	mu       sync.Mutex
	resolveN int
	result   championdata.Resolution
	err      error
}

func (f *fakeResolver) Resolve(_ context.Context, _, _ string, _ int) (championdata.Resolution, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.resolveN++
	return f.result, f.err
}

func (f *fakeResolver) resolveCalls() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.resolveN
}

// fakeState records every write the poller makes, standing in for
// state.Manager without a real state package dependency in tests.
type fakeState struct {
	mu sync.Mutex

	championCalls int
	lastChampID   string
	lastSkinID    int

	statsCalls int
	lastKills  int
	lastLevel  int

	startCalls int
	lastStart  int64
}

func (f *fakeState) UpdateChampion(championID, championName, skinName, chromaName string, skinID int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.championCalls++
	f.lastChampID = championID
	f.lastSkinID = skinID
}

func (f *fakeState) UpdateInGameStats(kills, deaths, assists, cs, level, gold int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.statsCalls++
	f.lastKills = kills
	f.lastLevel = level
}

func (f *fakeState) UpdateGameStart(start int64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.startCalls++
	f.lastStart = start
}

func (f *fakeState) snapshot() fakeState {
	f.mu.Lock()
	defer f.mu.Unlock()
	return fakeState{
		championCalls: f.championCalls, lastChampID: f.lastChampID, lastSkinID: f.lastSkinID,
		statsCalls: f.statsCalls, lastKills: f.lastKills, lastLevel: f.lastLevel,
		startCalls: f.startCalls, lastStart: f.lastStart,
	}
}

func waitFor(t *testing.T, timeout time.Duration, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("condition not met within %s", timeout)
}

func normalGameDoer() *fakeDoer {
	d := newFakeDoer()
	d.responses[baseURL+"/liveclientdata/activeplayer"] = `{"riotId":"Foo#EUW","level":0}`
	d.responses[baseURL+"/liveclientdata/allgamedata"] = `{"allPlayers":[
		{"riotId":"Foo#EUW","rawChampionName":"game_character_displayname_Chogath","rawSkinName":"game_character_skin_displayname_Chogath_1","skinID":1}
	]}`
	d.responses[baseURL+"/liveclientdata/playerscores?riotId=Foo%23EUW"] = `{"kills":1,"deaths":2,"assists":3,"creepScore":40}`
	d.responses[baseURL+"/liveclientdata/gamestats"] = `{"gameTime":100}`
	return d
}

func TestPoller_NormalGame_ResolvesChampionOnceAndPollsStatsEveryTick(t *testing.T) {
	d := normalGameDoer()
	resolver := &fakeResolver{result: championdata.Resolution{ID: "Chogath", Name: "Cho'Gath", BaseSkinID: 1}}
	state := &fakeState{}
	p := NewPoller(NewClient(d), resolver, state, storeWithInterval(testInterval), zerolog.Nop())

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	go p.Run(ctx, types.GameModeClassic)

	waitFor(t, testTimeout, func() bool { return state.snapshot().statsCalls >= 3 })

	snap := state.snapshot()
	if snap.championCalls != 1 {
		t.Errorf("championCalls = %d, want 1 (resolved once, cached)", snap.championCalls)
	}
	if snap.lastChampID != "Chogath" || snap.lastSkinID != 1 {
		t.Errorf("last champion write = %q/%d, want Chogath/1", snap.lastChampID, snap.lastSkinID)
	}
	if snap.lastKills != 1 {
		t.Errorf("lastKills = %d, want 1", snap.lastKills)
	}
	if resolver.resolveCalls() != 1 {
		t.Errorf("resolveCalls = %d, want 1 (not re-resolved every tick)", resolver.resolveCalls())
	}
	if snap.startCalls == 0 {
		t.Error("expected UpdateGameStart to be called")
	}
}

func TestPoller_NormalGame_LogsResolvedChampionSkinChroma(t *testing.T) {
	d := normalGameDoer()
	resolver := &fakeResolver{result: championdata.Resolution{
		ID: "Chogath", Name: "Cho'Gath", BaseSkinID: 1, SkinName: "Battlecast Cho'Gath", ChromaName: "Emerald",
	}}
	state := &fakeState{}

	var logBuf bytes.Buffer
	logger := zerolog.New(&logBuf)
	p := NewPoller(NewClient(d), resolver, state, storeWithInterval(testInterval), logger)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	go p.Run(ctx, types.GameModeClassic)

	waitFor(t, testTimeout, func() bool { return state.snapshot().championCalls == 1 })

	logged := logBuf.String()
	for _, want := range []string{"Cho'Gath", "Battlecast Cho'Gath", "Emerald", "resolved"} {
		if !strings.Contains(logged, want) {
			t.Errorf("log output missing %q; got: %s", want, logged)
		}
	}
}

func TestPoller_NormalGame_LoadingScreenSkipsChampionWriteButStillPollsStats(t *testing.T) {
	d := normalGameDoer()
	resolver := &fakeResolver{err: championdata.ErrUnresolved}
	state := &fakeState{}
	p := NewPoller(NewClient(d), resolver, state, storeWithInterval(testInterval), zerolog.Nop())

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	go p.Run(ctx, types.GameModeClassic)

	waitFor(t, testTimeout, func() bool { return state.snapshot().statsCalls >= 2 })

	if got := state.snapshot().championCalls; got != 0 {
		t.Errorf("championCalls = %d, want 0 while champion data is unresolved", got)
	}
}

func TestPoller_NormalGame_TransientScoreFailureSkipsThatTickWithoutClearing(t *testing.T) {
	d := normalGameDoer()
	delete(d.responses, baseURL+"/liveclientdata/playerscores?riotId=Foo%23EUW")
	resolver := &fakeResolver{result: championdata.Resolution{ID: "Chogath", Name: "Cho'Gath"}}
	state := &fakeState{}
	p := NewPoller(NewClient(d), resolver, state, storeWithInterval(testInterval), zerolog.Nop())

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	go p.Run(ctx, types.GameModeClassic)

	time.Sleep(20 * testInterval)

	if got := state.snapshot().statsCalls; got != 0 {
		t.Errorf("statsCalls = %d, want 0: a failed playerscores poll must not write zeroed stats", got)
	}
	// Champion resolution uses a different endpoint and should still succeed.
	waitFor(t, testTimeout, func() bool { return state.snapshot().championCalls == 1 })
}

func TestPoller_TFT_PollsLevelOnlyAndNeverTouchesChampionData(t *testing.T) {
	d := newFakeDoer()
	d.responses[baseURL+"/liveclientdata/activeplayer"] = `{"riotId":"Foo#EUW","level":5}`
	d.responses[baseURL+"/liveclientdata/gamestats"] = `{"gameTime":42}`
	resolver := &fakeResolver{result: championdata.Resolution{ID: "should-not-be-used"}}
	state := &fakeState{}
	p := NewPoller(NewClient(d), resolver, state, storeWithInterval(testInterval), zerolog.Nop())

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	go p.Run(ctx, types.GameModeTFT)

	waitFor(t, testTimeout, func() bool { return state.snapshot().statsCalls >= 2 })

	snap := state.snapshot()
	if snap.lastLevel != 5 {
		t.Errorf("lastLevel = %d, want 5", snap.lastLevel)
	}
	if snap.championCalls != 0 {
		t.Errorf("championCalls = %d, want 0 for TFT", snap.championCalls)
	}
	if resolver.resolveCalls() != 0 {
		t.Errorf("resolveCalls = %d, want 0 for TFT", resolver.resolveCalls())
	}
	if snap.startCalls == 0 {
		t.Error("expected UpdateGameStart to be called for TFT too")
	}
}

func TestPoller_TFT_MissingGameStatsStillUpdatesLevel(t *testing.T) {
	// gamestats isn't stubbed, so it 404s; level must still be written.
	d := newFakeDoer()
	d.responses[baseURL+"/liveclientdata/activeplayer"] = `{"riotId":"Foo#EUW","level":9}`
	resolver := &fakeResolver{}
	state := &fakeState{}
	p := NewPoller(NewClient(d), resolver, state, storeWithInterval(testInterval), zerolog.Nop())

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	go p.Run(ctx, types.GameModeTFT)

	waitFor(t, testTimeout, func() bool { return state.snapshot().statsCalls >= 1 })

	if got := state.snapshot().lastLevel; got != 9 {
		t.Errorf("lastLevel = %d, want 9 even though gamestats failed", got)
	}
}

func TestPoller_PicksUpStatsPollingIntervalChangeMidRun(t *testing.T) {
	t.Setenv("APPDATA", t.TempDir())
	d := normalGameDoer()
	resolver := &fakeResolver{result: championdata.Resolution{ID: "Chogath"}}
	state := &fakeState{}

	// Start well above the valid floor so only the initial tick fires soon.
	store := config.NewStore(&config.Config{StatsPollingInterval: 20000})
	p := NewPoller(NewClient(d), resolver, state, store, zerolog.Nop())

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	go p.Run(ctx, types.GameModeClassic)

	waitFor(t, testTimeout, func() bool { return state.snapshot().statsCalls == 1 })

	next := config.DefaultConfig()
	next.StatsPollingInterval = config.MinStatsPollingInterval // 1s, the fastest allowed
	if err := store.Apply(*next); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	waitFor(t, 5*time.Second, func() bool { return state.snapshot().statsCalls >= 2 })
}

func TestPoller_RunReturnsPromptlyOnContextCancel(t *testing.T) {
	d := normalGameDoer()
	resolver := &fakeResolver{result: championdata.Resolution{ID: "Chogath"}}
	state := &fakeState{}
	p := NewPoller(NewClient(d), resolver, state, storeWithInterval(time.Hour), zerolog.Nop()) // long interval; only the initial tick should fire

	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan struct{})
	go func() {
		p.Run(ctx, types.GameModeClassic)
		close(done)
	}()

	waitFor(t, testTimeout, func() bool { return state.snapshot().statsCalls >= 1 })
	cancel()

	select {
	case <-done:
	case <-time.After(testTimeout):
		t.Fatal("Run did not return after ctx cancellation")
	}
}
