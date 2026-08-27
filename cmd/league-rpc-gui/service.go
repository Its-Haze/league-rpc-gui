package main

import (
	"context"

	"github.com/its-haze/league-rpc/internal/app"
	"github.com/its-haze/league-rpc/internal/config"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// configChangedEvent is emitted to the frontend whenever the live config
// changes, carrying the new settings so screens don't have to re-fetch.
const configChangedEvent = "settings:changed"

// guiService adapts internal/app (plain Go, Wails-free) to a Wails service.
type guiService struct {
	app *app.App
}

func newGUIService(a *app.App) *guiService {
	return &guiService{app: a}
}

// GetSettings returns the current settings tree.
func (s *guiService) GetSettings() config.Config {
	return s.app.GetSettings()
}

// ApplySettings validates and persists cfg, swapping it in live. A returned
// error means nothing changed and the message is meant for the user.
func (s *guiService) ApplySettings(cfg config.Config) error {
	return s.app.ApplySettings(cfg)
}

// GetPresets returns the named Discord Application ID choices.
func (s *guiService) GetPresets() map[string]string {
	return s.app.GetPresets()
}

// publishConfigChanges forwards every config.Store update to the frontend as
// a Wails event until ctx is canceled.
func (s *guiService) publishConfigChanges(ctx context.Context, wailsApp *application.App) {
	updates := s.app.SubscribeSettings()
	for {
		select {
		case <-ctx.Done():
			return
		case cfg, ok := <-updates:
			if !ok {
				return
			}
			wailsApp.Event.Emit(configChangedEvent, cfg)
		}
	}
}
