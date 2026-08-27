package daemon

import (
	"context"
	"sync"
	"time"

	"github.com/its-haze/league-rpc/internal/discord"
	"github.com/its-haze/league-rpc/internal/state"
	"github.com/its-haze/league-rpc/pkg/types"
	"github.com/rs/zerolog"
)

// runner is anything Daemon can start and supervise for the duration of a
// run. *Supervisor and *LCUSupervisor both satisfy this.
type runner interface {
	Run(ctx context.Context)
}

// discordRunner also reports whether Discord is connected, so Daemon can
// resend presence the moment it reconnects. *Supervisor satisfies this.
type discordRunner interface {
	runner
	Connected() bool
}

// lcuRunner also reports LCU connection and League-process state, the
// inputs Daemon needs to pick real presence, placeholder, or cleared.
type lcuRunner interface {
	runner
	Connected() bool
	LeagueProcessDetected() bool
}

// liveGamePoller polls Live Client Data for champion/skin/KDA/timer detail
// while GameFlowInProgress is active. *livegame.Poller satisfies this.
type liveGamePoller interface {
	Run(ctx context.Context, gameMode types.GameMode)
}

// presenceMode is what Daemon currently shows on Discord. See ADR-0002.
type presenceMode int

const (
	modeUnknown presenceMode = iota
	modeConnected
	modePlaceholder
	modeCleared
)

// Daemon owns the Discord and LCU Connection Supervisors and drives Discord
// presence off their state. Run is the single seam here. See ADR-0002.
type Daemon struct {
	discord  discordRunner
	lcu      lcuRunner
	updater  *discord.Updater
	state    *state.Manager
	liveGame liveGamePoller
	logger   zerolog.Logger

	presencePollInterval time.Duration
	placeholderInterval  time.Duration
}

// New builds a Daemon that drives presence from stateMgr through updater.
func New(
	discordSup discordRunner,
	lcuSup lcuRunner,
	updater *discord.Updater,
	stateMgr *state.Manager,
	liveGame liveGamePoller,
	logger zerolog.Logger,
	presencePollInterval, placeholderInterval time.Duration,
) *Daemon {
	return &Daemon{
		discord:              discordSup,
		lcu:                  lcuSup,
		updater:              updater,
		state:                stateMgr,
		liveGame:             liveGame,
		logger:               logger,
		presencePollInterval: presencePollInterval,
		placeholderInterval:  placeholderInterval,
	}
}

// Run starts both supervisors and the presence loop, and blocks until ctx
// is canceled; a connection failure on either side never returns early.
func (d *Daemon) Run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		d.discord.Run(ctx)
	}()
	go func() {
		defer wg.Done()
		d.lcu.Run(ctx)
	}()
	go func() {
		defer wg.Done()
		d.updater.Run(ctx)
	}()
	go func() {
		defer wg.Done()
		d.presenceLoop(ctx)
	}()

	wg.Wait()
}

// presenceLoop shows real presence once the LCU connects, a rotating
// placeholder while League is launching, or clears presence. See ADR-0002.
func (d *Daemon) presenceLoop(ctx context.Context) {
	ticker := time.NewTicker(d.presencePollInterval)
	defer ticker.Stop()

	cfgUpdates := d.updater.ConfigChanges()

	mode := modeUnknown
	discordConnected := false
	waitingForDiscordLogged := false

	var placeholderTicker *time.Ticker
	var placeholderC <-chan time.Time
	var placeholderStart int64

	stopPlaceholder := func() {
		if placeholderTicker != nil {
			placeholderTicker.Stop()
			placeholderTicker = nil
			placeholderC = nil
		}
	}
	defer stopPlaceholder()

	sendPlaceholder := func() {
		// Skip while Discord isn't connected yet; the reconnect handler below resends once it is.
		if !d.discord.Connected() {
			return
		}
		d.updater.UpdatePlaceholder(discord.BuildLaunchingPresence(placeholderStart))
	}

	var liveGameCancel context.CancelFunc
	liveGameActive := false

	startLiveGame := func(gameMode types.GameMode) {
		if liveGameActive {
			return
		}
		liveGameActive = true
		// The previous match's in-game data is already cleared by
		var gameCtx context.Context
		gameCtx, liveGameCancel = context.WithCancel(ctx)
		go d.liveGame.Run(gameCtx, gameMode)
	}
	stopLiveGame := func() {
		if !liveGameActive {
			return
		}
		liveGameActive = false
		liveGameCancel()
		liveGameCancel = nil
	}
	defer stopLiveGame()

	for {
		select {
		case <-ctx.Done():
			return

		case st, ok := <-d.state.Updates():
			if !ok {
				return
			}
			if mode == modeConnected {
				d.updater.DelayUpdate(st)
			}

		case <-cfgUpdates:
			// Display settings may have changed; reflect them now instead of
			// waiting for the next real state change or poll tick.
			if mode == modeConnected {
				d.updater.ImmediateUpdate(d.state.Get())
			}

		case <-placeholderC:
			sendPlaceholder()

		case <-ticker.C:
			// Discord can connect after the LCU already has; resend the
			// current mode's presence instead of waiting for a state change.
			nowConnected := d.discord.Connected()
			if nowConnected && !discordConnected {
				switch mode {
				case modeConnected:
					d.updater.ImmediateUpdate(d.state.Get())
				case modePlaceholder:
					sendPlaceholder()
				}
			}
			discordConnected = nowConnected

			// League is up but Discord isn't reachable yet: log it once per
			// edge, not on every poll tick, so it isn't spammy.
			if d.lcu.LeagueProcessDetected() && !nowConnected {
				if !waitingForDiscordLogged {
					d.logger.Info().Msg("League is running, waiting for Discord to be reachable")
					waitingForDiscordLogged = true
				}
			} else {
				waitingForDiscordLogged = false
			}

			switch {
			case d.lcu.Connected():
				st := d.state.Get()
				if mode != modeConnected {
					stopPlaceholder()
					// Skip while Discord isn't connected yet; the reconnect handler above resends once it is.
					if d.discord.Connected() {
						d.updater.ImmediateUpdate(st)
					}
					mode = modeConnected
				}

				if st.GameFlowPhase == types.GameFlowInProgress {
					startLiveGame(st.GameMode)
				} else {
					stopLiveGame()
				}

			case d.lcu.LeagueProcessDetected():
				if mode != modePlaceholder {
					mode = modePlaceholder
					placeholderStart = time.Now().Unix()
					sendPlaceholder()
					placeholderTicker = time.NewTicker(d.placeholderInterval)
					placeholderC = placeholderTicker.C
				}

			default:
				if mode != modeCleared {
					stopPlaceholder()
					d.updater.ClearPresence()
					mode = modeCleared
				}
			}

			// The live game poller only ever runs while mode == modeConnected
			if mode != modeConnected {
				stopLiveGame()
			}
		}
	}
}
