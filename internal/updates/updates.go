// Package updates checks GitHub Releases and, on an explicit user action,
// downloads, verifies, and swaps the binary. It never downloads on a timer.
package updates

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
	wupdater "github.com/wailsapp/wails/v3/pkg/updater"
)

// engine is the slice of the Wails updater Coordinator drives. *wupdater.Updater
// satisfies it; tests supply a fake.
type engine interface {
	Check(ctx context.Context) (*wupdater.Release, error)
	DownloadAndInstall(ctx context.Context) error
	Restart(ctx context.Context) error
}

// Status is what the GUI renders for the App Update state. Available flips true
// only while a check has found a strictly newer release.
type Status struct {
	Available bool   `json:"available"`
	Version   string `json:"version"`
	Notes     string `json:"notes"`
	LastError string `json:"last_error,omitempty"`
}

// Coordinator owns the check cadence and the download/restart actions; a dev
// build disables it entirely, since there is no released version to compare.
type Coordinator struct {
	eng       engine
	changelog *changelogFetcher
	log       zerolog.Logger
	dev       bool
	interval  time.Duration

	mu       sync.Mutex
	status   Status
	onChange func(Status)
}

// New builds a Coordinator. dev=true disables every check/download action.
// doer backs the changelog fetch and is injected so tests never hit GitHub.
func New(eng engine, doer HTTPDoer, dev bool, log zerolog.Logger) *Coordinator {
	return &Coordinator{
		eng:       eng,
		changelog: newChangelogFetcher(doer, RepoSlug),
		log:       log,
		dev:       dev,
		interval:  CheckInterval,
	}
}

// OnChange registers the callback fired whenever the status changes. The GUI
// adapter forwards it to a frontend event.
func (c *Coordinator) OnChange(fn func(Status)) {
	c.mu.Lock()
	c.onChange = fn
	c.mu.Unlock()
}

// Status returns the last known App Update status.
func (c *Coordinator) Status() Status {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.status
}

// Run does a launch check, then re-checks every CheckInterval until ctx is
// canceled. A dev build returns immediately.
func (c *Coordinator) Run(ctx context.Context) {
	if c.dev {
		c.log.Info().Msg("dev build: App Update checks disabled")
		return
	}
	c.check(ctx)

	t := time.NewTicker(c.interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			c.check(ctx)
		}
	}
}

// Check runs one check now and reports the resulting status. Used by the
// manual "Check for updates" action.
func (c *Coordinator) Check(ctx context.Context) (Status, error) {
	if c.dev {
		return c.Status(), nil
	}
	return c.check(ctx)
}

// check performs one provider round trip; an error keeps any prior available
// release (recording the message), and a clean "no update" clears it.
func (c *Coordinator) check(ctx context.Context) (Status, error) {
	rel, err := c.eng.Check(ctx)

	c.mu.Lock()
	switch {
	case err != nil:
		c.status.LastError = err.Error()
		c.log.Warn().Err(err).Msg("App Update check failed")
	case rel == nil:
		c.status = Status{}
	default:
		c.status = Status{Available: true, Version: rel.Version, Notes: rel.Notes}
		c.log.Info().Str("version", rel.Version).Msg("App Update available")
	}
	s := c.status
	cb := c.onChange
	c.mu.Unlock()

	if cb != nil {
		cb(s)
	}
	return s, err
}

// Download re-checks (through check, so status and onChange stay current)
// then downloads, verifies against SHA256SUMS, and swaps. Never restarts.
func (c *Coordinator) Download(ctx context.Context) error {
	if c.dev {
		return errors.New("updates: a dev build cannot self-update")
	}
	s, err := c.check(ctx)
	if err != nil {
		return fmt.Errorf("updates: re-check before download: %w", err)
	}
	if !s.Available {
		return errors.New("updates: no update available")
	}
	if err := c.eng.DownloadAndInstall(ctx); err != nil {
		return fmt.Errorf("updates: download and install: %w", err)
	}
	return nil
}

// Restart relaunches into the freshly swapped binary. Only valid after a
// successful Download.
func (c *Coordinator) Restart(ctx context.Context) error {
	if c.dev {
		return errors.New("updates: a dev build cannot self-update")
	}
	return c.eng.Restart(ctx)
}

// Changelog returns the latest release's notes as Markdown, or the literal
// "changelog unavailable" when GitHub can't be reached.
func (c *Coordinator) Changelog(ctx context.Context) string {
	body, err := c.changelog.fetch(ctx)
	if err != nil {
		c.log.Debug().Err(err).Msg("changelog fetch failed")
		return ChangelogUnavailable
	}
	return body
}
