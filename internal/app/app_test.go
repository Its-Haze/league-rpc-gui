package app

import (
	"reflect"
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

func TestApp_RenderTemplatePreview_UsesSampleDataAndDefaults(t *testing.T) {
	a := New(config.NewStore(config.DefaultConfig()), &fakePauser{})

	dt, st := config.DefaultConfig().Presence.Templates["in-game"].Details, config.DefaultConfig().Presence.Templates["in-game"].State
	got, err := a.RenderTemplatePreview("in-game", config.TemplatePair{Details: dt, State: st}, nil)
	if err != nil {
		t.Fatalf("RenderTemplatePreview: %v", err)
	}
	if got.Details != "Ranked Solo/Duo" || got.State != "In Game · 3/2/5 · 120cs" {
		t.Fatalf("preview = %+v, want the sample-data render", got)
	}
	if len(got.Warnings) != 0 {
		t.Fatalf("warnings = %v, want none", got.Warnings)
	}
}

func TestApp_RenderTemplatePreview_ReportsUnknownTokens(t *testing.T) {
	a := New(config.NewStore(config.DefaultConfig()), &fakePauser{})

	got, err := a.RenderTemplatePreview("in-client",
		config.TemplatePair{Details: "{availability} {foo}", State: "{bar}"},
		map[string]string{"availability": "Online"})
	if err != nil {
		t.Fatalf("RenderTemplatePreview: %v", err)
	}
	if got.Details != "Online {foo}" || got.State != "{bar}" {
		t.Fatalf("preview = %+v, want unknown tokens left literal", got)
	}
	if want := []string{"unknown token {bar}", "unknown token {foo}"}; !reflect.DeepEqual(got.Warnings, want) {
		t.Fatalf("warnings = %v, want %v", got.Warnings, want)
	}
}

func TestApp_RenderTemplatePreview_RejectsUnknownContext(t *testing.T) {
	a := New(config.NewStore(config.DefaultConfig()), &fakePauser{})
	if _, err := a.RenderTemplatePreview("bogus", config.TemplatePair{}, nil); err == nil {
		t.Fatal("accepted an unknown presence context")
	}
}

// A blank line in the preview must resolve to the same default the daemon uses,
// so the settings screen can't show something Discord won't.
func TestApp_RenderTemplatePreview_BlankLineMatchesRuntimeDefault(t *testing.T) {
	a := New(config.NewStore(config.DefaultConfig()), &fakePauser{})

	got, err := a.RenderTemplatePreview("in-game", config.TemplatePair{Details: "", State: "custom state"}, nil)
	if err != nil {
		t.Fatalf("RenderTemplatePreview: %v", err)
	}
	if got.Details != "Ranked Solo/Duo" {
		t.Fatalf("blank Details previewed as %q, want the built-in default", got.Details)
	}
	if got.State != "custom state" {
		t.Fatalf("State = %q, want the override", got.State)
	}
}

// A non-nil but empty sample map (what the JS boundary produces from {}) must
// still fall back to representative data, not render every token empty.
func TestApp_RenderTemplatePreview_EmptyMapUsesSampleData(t *testing.T) {
	a := New(config.NewStore(config.DefaultConfig()), &fakePauser{})

	got, err := a.RenderTemplatePreview("champ-select", config.TemplatePair{}, map[string]string{})
	if err != nil {
		t.Fatalf("RenderTemplatePreview: %v", err)
	}
	if got.Details == "" {
		t.Fatal("empty sample map produced a blank preview")
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
