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
	// pauseHook, when set, handles pause changes so the tray checkbox stays
	// in sync. Nil until the tray is wired; falls back to app then.
	pauseHook func(bool)
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

// RenderTemplatePreview renders a presence template pair for ctx against sample
// data so the settings screen can preview an edit before it is saved.
func (s *guiService) RenderTemplatePreview(ctx string, tmpl config.TemplatePair, sample map[string]string) (app.TemplatePreview, error) {
	return s.app.RenderTemplatePreview(ctx, tmpl, sample)
}

// SetPaused toggles the runtime pause flag from the frontend.
func (s *guiService) SetPaused(paused bool) {
	if s.pauseHook != nil {
		s.pauseHook(paused)
		return
	}
	s.app.SetPaused(paused)
}

// IsPaused reports the runtime pause flag to the frontend.
func (s *guiService) IsPaused() bool {
	return s.app.IsPaused()
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
