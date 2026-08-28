// Package app is the seam the GUI binds to. Its methods are plain Go so they
// can be exposed as Wails bindings without this package importing Wails.
package app

import (
	"fmt"
	"sort"

	"github.com/its-haze/league-rpc/internal/config"
	"github.com/its-haze/league-rpc/internal/presence/template"
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
}

// New builds an App over store and the daemon's pause control.
func New(store *config.Store, pauser Pauser) *App {
	return &App{store: store, pauser: pauser}
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
