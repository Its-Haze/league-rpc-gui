package config

import (
	"encoding/json"
	"os"
	"testing"
)

const flatV1Fixture = `{
  "discord_app_id": "999",
  "show_stats": false,
  "show_rank": true,
  "show_emojis": false,
  "show_in_client": true,
  "auto_launch_league": true,
  "league_path": "C:\\Games\\League",
  "update_interval": 2500,
  "stats_polling_interval": 4000,
  "debug_mode": true
}`

func TestLoad_MigratesFlatFileAndWritesItBack(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("APPDATA", dir)

	path, _ := GetConfigPath()
	if err := os.MkdirAll(dirOf(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(flatV1Fixture), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if got.SchemaVersion != CurrentSchemaVersion {
		t.Errorf("SchemaVersion = %d, want %d", got.SchemaVersion, CurrentSchemaVersion)
	}
	if got.DiscordAppID != "999" {
		t.Errorf("DiscordAppID = %q", got.DiscordAppID)
	}
	if got.Display.Default.ShowStats {
		t.Error("show_stats=false did not map to Display.Default.ShowStats")
	}
	if !got.Display.Default.ShowRank {
		t.Error("show_rank=true did not map to Display.Default.ShowRank")
	}
	if got.Presence.ShowEmojis {
		t.Error("show_emojis=false did not map to Presence.ShowEmojis")
	}
	if !got.Behavior.AutoLaunchLeague || got.Behavior.LeaguePath != `C:\Games\League` {
		t.Errorf("league fields did not migrate: %+v", got.Behavior)
	}
	if got.Advanced.UpdateInterval != 2500 || got.Advanced.StatsPollingInterval != 4000 || !got.Advanced.DebugMode {
		t.Errorf("advanced fields did not migrate: %+v", got.Advanced)
	}
	if got.Theme != ThemeSystem {
		t.Errorf("Theme = %q, want default %q", got.Theme, ThemeSystem)
	}

	// The upgraded file is on disk: reloading takes the non-migration path.
	raw, _ := os.ReadFile(path)
	if schemaVersionOf(raw) != CurrentSchemaVersion {
		t.Fatalf("file was not written back with a schema version: %s", raw)
	}
}

func TestLoad_V2FileRoundTrips(t *testing.T) {
	t.Setenv("APPDATA", t.TempDir())

	orig := DefaultConfig()
	orig.Theme = ThemeDark
	orig.Display.Default.ShowRank = false
	orig.Display.Modes["ARENA"] = ModeOverride{}
	if err := Save(orig); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	a, _ := json.Marshal(orig)
	b, _ := json.Marshal(got)
	if string(a) != string(b) {
		t.Fatalf("v2 config did not round-trip:\n%s\n%s", a, b)
	}
}

func TestLoad_BrokenFileClampsToValid(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("APPDATA", dir)

	path, _ := GetConfigPath()
	if err := os.MkdirAll(dirOf(path), 0o755); err != nil {
		t.Fatal(err)
	}
	broken := `{"schema_version": 2, "discord_app_id": "", "theme": "chartreuse",
	  "advanced": {"update_interval": 1, "stats_polling_interval": 9999999}}`
	if err := os.WriteFile(path, []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("Load returned an invalid config: %v", err)
	}
}

func dirOf(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' || p[i] == '\\' {
			return p[:i]
		}
	}
	return "."
}
