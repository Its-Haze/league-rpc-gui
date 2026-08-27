package config

import "testing"

func TestValidate_PureAndReportsEveryProblem(t *testing.T) {
	c := Config{DiscordAppID: "", UpdateInterval: 10, StatsPollingInterval: 999999}
	before := c

	err := c.Validate()
	if err == nil {
		t.Fatal("Validate accepted an invalid config")
	}
	if c != before {
		t.Fatalf("Validate mutated the config: %+v -> %+v", before, c)
	}
	for _, want := range []string{"discord_app_id", "update_interval", "stats_polling_interval"} {
		if !contains(err.Error(), want) {
			t.Errorf("error %q missing mention of %q", err, want)
		}
	}
}

func TestValidate_AcceptsDefault(t *testing.T) {
	if err := DefaultConfig().Validate(); err != nil {
		t.Fatalf("DefaultConfig failed validation: %v", err)
	}
}

func TestClamp_RepairsOutOfBounds(t *testing.T) {
	c := Config{DiscordAppID: "", UpdateInterval: 50, StatsPollingInterval: 100}
	c.clamp()

	if c.DiscordAppID == "" {
		t.Error("clamp left DiscordAppID empty")
	}
	if c.UpdateInterval != DefaultConfig().UpdateInterval {
		t.Errorf("UpdateInterval = %d, want default", c.UpdateInterval)
	}
	if c.StatsPollingInterval != MinStatsPollingInterval {
		t.Errorf("StatsPollingInterval = %d, want %d", c.StatsPollingInterval, MinStatsPollingInterval)
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("config still invalid after clamp: %v", err)
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
