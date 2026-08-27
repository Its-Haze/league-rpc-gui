package app

import (
	"testing"

	"github.com/its-haze/league-rpc/internal/config"
)

func TestApp_ApplySettingsPersistsAndGetReflects(t *testing.T) {
	t.Setenv("APPDATA", t.TempDir())
	store := config.NewStore(config.DefaultConfig())
	a := New(store)

	cfg := a.GetSettings()
	cfg.ShowStats = false
	cfg.UpdateInterval = 2500

	if err := a.ApplySettings(cfg); err != nil {
		t.Fatalf("ApplySettings failed: %v", err)
	}

	got := a.GetSettings()
	if got.ShowStats || got.UpdateInterval != 2500 {
		t.Fatalf("GetSettings did not reflect the change: %+v", got)
	}
}

func TestApp_ApplySettingsRejectsInvalid(t *testing.T) {
	t.Setenv("APPDATA", t.TempDir())
	a := New(config.NewStore(config.DefaultConfig()))

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
	a := New(config.NewStore(config.DefaultConfig()))
	if len(a.GetPresets()) == 0 {
		t.Fatal("GetPresets returned nothing")
	}
}
