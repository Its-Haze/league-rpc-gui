// Package app is the seam the GUI binds to. Its methods are plain Go so they
// can be exposed as Wails bindings without this package importing Wails.
package app

import (
	"context"
	"fmt"
	"sort"

	"github.com/its-haze/league-rpc/internal/config"
	"github.com/its-haze/league-rpc/internal/presence/template"
	"github.com/its-haze/league-rpc/internal/state"
)

// Pauser is the runtime pause control the daemon exposes. Kept as a local
// interface so this package stays free of a daemon import.
type Pauser interface {
	SetPaused(bool)
	IsPaused() bool
}

// App exposes the settings surface to the frontend. Everything runtime-visible
// goes through the config.Store it holds; the daemon reads the same store.
type App struct {
	store  *config.Store
	pauser Pauser

	status *statusBridge
	tester TestPresenter
}

// Option configures optional App wiring the settings surface does not need.
type Option func(*App)

// WithStatus wires the status bridge (conns, probe, and state changes on
func WithStatus(conns Connections, probe PresenceProbe, states <-chan *state.State, tester TestPresenter) Option {
	return func(a *App) {
		a.status = newStatusBridge(conns, probe, a.pauser, states)
		a.tester = tester
	}
}

// New builds an App over store and the daemon's pause control.
func New(store *config.Store, pauser Pauser, opts ...Option) *App {
	a := &App{store: store, pauser: pauser}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// GetStatus returns the current status snapshot. Zero value until WithStatus
// is wired.
func (a *App) GetStatus() StatusSnapshot {
	if a.status == nil {
		return StatusSnapshot{}
	}
	return a.status.snapshot()
}

// OnStatusChange registers the callback fired whenever the status snapshot
// changes. The GUI adapter forwards it to a frontend event.
func (a *App) OnStatusChange(fn func(StatusSnapshot)) {
	if a.status == nil {
		return
	}
	a.status.setOnChange(fn)
}

// RunStatus drives the status bridge until ctx is canceled. Without WithStatus
// it just blocks until then. The GUI adapter runs this in a goroutine.
func (a *App) RunStatus(ctx context.Context) {
	if a.status == nil {
		<-ctx.Done()
		return
	}
	a.status.run(ctx)
}

// TestPresence shows a fixed sample presence in Discord for a few seconds so
// the user can confirm their settings without launching League.
func (a *App) TestPresence() {
	if a.tester == nil {
		return
	}
	a.tester.TestPresence()
}

// GetSettings returns the current settings as a value copy for the frontend.
func (a *App) GetSettings() config.Config {
	return *a.store.Load()
}

// ApplySettings validates cfg, persists it, and swaps it in live. A returned
// error means nothing changed and the message is meant for the user.
func (a *App) ApplySettings(cfg config.Config) error {
	return a.store.Apply(cfg)
}

// GetPresets returns the named Discord Application ID choices for the picker.
func (a *App) GetPresets() map[string]string {
	return config.DiscordAppIDPresets()
}

// SetPaused toggles the runtime pause flag. Paused clears Discord presence
// immediately; it never persists to Config and resets to unpaused on restart.
func (a *App) SetPaused(paused bool) {
	a.pauser.SetPaused(paused)
}

// IsPaused reports the runtime pause flag.
func (a *App) IsPaused() bool {
	return a.pauser.IsPaused()
}

// SubscribeSettings returns a channel that receives the new settings on every
// successful ApplySettings. The GUI adapter forwards these to the frontend.
func (a *App) SubscribeSettings() <-chan *config.Config {
	return a.store.Subscribe()
}

// TemplatePreview is the rendered result of one presence template pair against
// sample data, plus any unknown-token warnings the engine reported.
type TemplatePreview struct {
	Details  string   `json:"details"`
	State    string   `json:"state"`
	Warnings []string `json:"warnings"`
}

// RenderTemplatePreview renders tmpl for ctx the same way the daemon will: a
// blank line falls back to the default, an empty sample to the built-in sample.
func (a *App) RenderTemplatePreview(ctx string, tmpl config.TemplatePair, sample map[string]string) (TemplatePreview, error) {
	tctx := template.Context(ctx)
	if !template.IsContext(tctx) {
		return TemplatePreview{}, fmt.Errorf("unknown presence context %q", ctx)
	}
	if len(sample) == 0 {
		sample = template.SampleData(tctx)
	}

	details, state, unknown := template.RenderPair(tctx, tmpl.Details, tmpl.State, sample)

	var warnings []string
	for _, name := range unknown {
		warnings = append(warnings, fmt.Sprintf("unknown token {%s}", name))
	}
	sort.Strings(warnings)

	return TemplatePreview{Details: details, State: state, Warnings: warnings}, nil
}
