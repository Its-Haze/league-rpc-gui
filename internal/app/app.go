// Package app is the seam the GUI binds to. Its methods are plain Go so they
// can be exposed as Wails bindings without this package importing Wails.
package app

import "github.com/its-haze/league-rpc/internal/config"

// App exposes the settings surface to the frontend. Everything runtime-visible
// goes through the config.Store it holds; the daemon reads the same store.
type App struct {
	store *config.Store
}

// New builds an App over store.
func New(store *config.Store) *App {
	return &App{store: store}
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

// SubscribeSettings returns a channel that receives the new settings on every
// successful ApplySettings. The GUI adapter forwards these to the frontend.
func (a *App) SubscribeSettings() <-chan *config.Config {
	return a.store.Subscribe()
}
