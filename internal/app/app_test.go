package app

import (
	"testing"

	"github.com/its-haze/league-rpc/internal/config"
)

// fakePauser records the last SetPaused value and reports it from IsPaused.
type fakePauser struct{ paused bool }

func (f *fakePauser) SetPaused(p bool) { f.paused = p }
func (f *fakePauser) IsPaused() bool   { return f.paused }

func TestApp_ApplySettingsPersistsAndGetReflects(t *testing.T) {
	t.Setenv("APPDATA", t.TempDir())
	store := config.NewStore(config.DefaultConfig())
	a := New(store, &fakePauser{})

	cfg := a.GetSettings()
	cfg.Display.Default.ShowStats = false
	cfg.Advanced.UpdateInterval = 2500

	if err := a.ApplySettings(cfg); err != nil {
		t.Fatalf("ApplySettings failed: %v", err)
	}

	got := a.GetSettings()
	if got.Display.Default.ShowStats || got.Advanced.UpdateInterval != 2500 {
		t.Fatalf("GetSettings did not reflect the change: %+v", got)
	}
}

func TestApp_ApplySettingsRejectsInvalid(t *testing.T) {
	t.Setenv("APPDATA", t.TempDir())
	a := New(config.NewStore(config.DefaultConfig()), &fakePauser{})

	cfg := a.GetSettings()
	cfg.DiscordAppID = ""

	if err := a.ApplySettings(cfg); err == nil {
		t.Fatal("ApplySettings accepted an empty DiscordAppID")
	}
	if a.GetSettings().DiscordAppID == "" {
		t.Fatal("rejected ApplySettings still mutated live settings")
	}
}

func TestApp_GetPresetsNonEmpty(t *testing.T) {
	a := New(config.NewStore(config.DefaultConfig()), &fakePauser{})
	if len(a.GetPresets()) == 0 {
		t.Fatal("GetPresets returned nothing")
	}
}

func TestApp_PauseDelegatesToPauser(t *testing.T) {
	p := &fakePauser{}
	a := New(config.NewStore(config.DefaultConfig()), p)

	if a.IsPaused() {
		t.Fatal("expected unpaused initially")
	}
	a.SetPaused(true)
	if !p.paused || !a.IsPaused() {
		t.Fatal("SetPaused(true) did not reach the pauser")
	}
	a.SetPaused(false)
	if p.paused || a.IsPaused() {
		t.Fatal("SetPaused(false) did not reach the pauser")
	}
}
