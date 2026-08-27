package config

import (
	"encoding/json"
	"testing"
)

func TestValidate_ReportsEveryProblemAtOnce(t *testing.T) {
	c := DefaultConfig()
	c.DiscordAppID = ""
	c.Theme = "neon"
	c.Advanced.UpdateInterval = 10
	c.Advanced.StatsPollingInterval = 999999

	err := c.Validate()
	if err == nil {
		t.Fatal("Validate accepted an invalid config")
	}
	for _, want := range []string{"discord_app_id", "theme", "update_interval", "stats_polling_interval"} {
		if !contains(err.Error(), want) {
			t.Errorf("error %q missing mention of %q", err, want)
		}
	}
}

func TestValidate_DoesNotMutate(t *testing.T) {
	c := DefaultConfig()
	c.DiscordAppID = ""
	c.Advanced.UpdateInterval = 10

	before, _ := json.Marshal(c)
	_ = c.Validate()
	after, _ := json.Marshal(c)

	if string(before) != string(after) {
		t.Fatalf("Validate mutated the config:\n%s\n%s", before, after)
	}
}

func TestValidate_AcceptsDefault(t *testing.T) {
	if err := DefaultConfig().Validate(); err != nil {
		t.Fatalf("DefaultConfig failed validation: %v", err)
	}
}

func TestClamp_RepairsOutOfBounds(t *testing.T) {
	c := &Config{
		Theme: "bogus",
		Advanced: AdvancedConfig{
			UpdateInterval:       50,
			StatsPollingInterval: 100,
		},
	}
	c.clamp()

	if c.DiscordAppID == "" {
		t.Error("clamp left DiscordAppID empty")
	}
	if c.Theme != ThemeSystem {
		t.Errorf("Theme = %q, want %q", c.Theme, ThemeSystem)
	}
	if c.Advanced.UpdateInterval != DefaultConfig().Advanced.UpdateInterval {
		t.Errorf("UpdateInterval = %d, want default", c.Advanced.UpdateInterval)
	}
	if c.Advanced.StatsPollingInterval != MinStatsPollingInterval {
		t.Errorf("StatsPollingInterval = %d, want %d", c.Advanced.StatsPollingInterval, MinStatsPollingInterval)
	}
	if c.Display.Modes == nil || c.Presence.Templates == nil {
		t.Error("clamp left nested maps nil")
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("config still invalid after clamp: %v", err)
	}
}

func TestDisplayConfig_Resolve(t *testing.T) {
	f := false
	d := DisplayConfig{
		Default: DisplayDefaults{ShowRank: true, ShowStats: true},
		Modes: map[string]ModeOverride{
			"ARENA": {ShowRank: &f},
		},
	}

	if got := d.Resolve("CLASSIC"); !got.ShowRank || !got.ShowStats {
		t.Errorf("unknown mode should inherit default, got %+v", got)
	}
	got := d.Resolve("ARENA")
	if got.ShowRank {
		t.Error("ARENA override of ShowRank=false was not applied")
	}
	if !got.ShowStats {
		t.Error("ARENA should still inherit ShowStats=true")
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
