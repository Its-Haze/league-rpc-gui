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

// statusChangedEvent carries a StatusSnapshot to the frontend whenever the
// connection state, phase, or last-sent presence changes.
const statusChangedEvent = "status:changed"

// updateChangedEvent carries an app.UpdateStatus to the frontend whenever the
// App Update coordinator's launch or periodic check finds something new.
const updateChangedEvent = "update:changed"

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

// GetStatus returns the current status snapshot for the frontend.
func (s *guiService) GetStatus() app.StatusSnapshot {
	return s.app.GetStatus()
}

// TestPresence shows a fixed sample presence in Discord for a few seconds.
func (s *guiService) TestPresence() {
	s.app.TestPresence()
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

// GetVersion returns the running build's version, or a dev placeholder.
func (s *guiService) GetVersion() string {
	return s.app.GetVersion()
}

// GetUpdateStatus returns the last known App Update status.
func (s *guiService) GetUpdateStatus() app.UpdateStatus {
	return s.app.GetUpdateStatus()
}

// CheckForUpdates runs the manual "Check for updates" action.
func (s *guiService) CheckForUpdates(ctx context.Context) (app.UpdateStatus, error) {
	return s.app.CheckForUpdates(ctx)
}

// StartUpdate downloads, verifies, and swaps the pending release.
func (s *guiService) StartUpdate(ctx context.Context) error {
	return s.app.StartUpdate(ctx)
}

// RestartForUpdate relaunches into the freshly swapped binary.
func (s *guiService) RestartForUpdate(ctx context.Context) error {
	return s.app.RestartForUpdate(ctx)
}

// GetChangelog returns the latest release's notes as Markdown, or a fixed
// placeholder when GitHub can't be reached.
func (s *guiService) GetChangelog(ctx context.Context) string {
	return s.app.GetChangelog(ctx)
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
