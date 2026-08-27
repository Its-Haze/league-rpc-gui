package config

import (
	"errors"
	"fmt"
)

// Config represents the application configuration
// All settings are stored in a JSON file and managed through the GUI
type Config struct {
	// Discord Settings
	DiscordAppID string `json:"discord_app_id"` // Discord Application ID

	// Display Settings
	ShowStats    bool `json:"show_stats"`     // Show KDA and Creep Score
	ShowRank     bool `json:"show_rank"`      // Show rank emblem and LP
	ShowEmojis   bool `json:"show_emojis"`    // Show online/away emoji indicators
	ShowInClient bool `json:"show_in_client"` // Show RPC when idle in client

	// League Settings
	AutoLaunchLeague bool   `json:"auto_launch_league"` // Auto-launch League on startup
	LeaguePath       string `json:"league_path"`        // Custom League installation path (optional)

	// Advanced Settings
	UpdateInterval       int  `json:"update_interval"`        // RPC update throttle in milliseconds
	StatsPollingInterval int  `json:"stats_polling_interval"` // In-game stats polling interval in milliseconds
	DebugMode            bool `json:"debug_mode"`             // Enable debug logging
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		DiscordAppID:         "1237146703111393281", // TEMP dev app ID; real one is 1194034071588851783
		ShowStats:            true,
		ShowRank:             true,
		ShowEmojis:           true,
		ShowInClient:         true,
		AutoLaunchLeague:     false,
		LeaguePath:           "",   // Auto-detect by default
		UpdateInterval:       1500, // 1.5 seconds
		StatsPollingInterval: 3000, // 3 seconds
		DebugMode:            false,
	}
}

// Bounds for the numeric settings. Shared by Validate and clamp so the
// reject path and the load-time repair path can't drift apart.
const (
	MinUpdateInterval       = 500
	MaxUpdateInterval       = 10000
	MinStatsPollingInterval = 1000
	MaxStatsPollingInterval = 30000
)

// Validate reports every problem with c without changing it. Store.Apply and
// Save call this and surface the error; the GUI shows it next to the field.
func (c *Config) Validate() error {
	var errs []error

	if c.DiscordAppID == "" {
		errs = append(errs, errors.New("discord_app_id must not be empty"))
	}
	if c.UpdateInterval < MinUpdateInterval || c.UpdateInterval > MaxUpdateInterval {
		errs = append(errs, fmt.Errorf("update_interval must be between %d and %d ms", MinUpdateInterval, MaxUpdateInterval))
	}
	if c.StatsPollingInterval < MinStatsPollingInterval || c.StatsPollingInterval > MaxStatsPollingInterval {
		errs = append(errs, fmt.Errorf("stats_polling_interval must be between %d and %d ms", MinStatsPollingInterval, MaxStatsPollingInterval))
	}

	return errors.Join(errs...)
}

// clamp forces c into valid bounds in place. Only the load path uses it, so a
// hand-broken or older config file on disk still boots with sane values.
func (c *Config) clamp() {
	if c.DiscordAppID == "" {
		c.DiscordAppID = DefaultConfig().DiscordAppID
	}
	if c.UpdateInterval < MinUpdateInterval || c.UpdateInterval > MaxUpdateInterval {
		c.UpdateInterval = DefaultConfig().UpdateInterval
	}
	if c.StatsPollingInterval < MinStatsPollingInterval {
		c.StatsPollingInterval = MinStatsPollingInterval
	}
	if c.StatsPollingInterval > MaxStatsPollingInterval {
		c.StatsPollingInterval = MaxStatsPollingInterval
	}
}

// DiscordAppIDPresets returns a map of Discord App ID presets
func DiscordAppIDPresets() map[string]string {
	return map[string]string{
		"League of Legends": "1194034071588851783",
		"League of Kittens": "1230607224296968303",
		"League of Linux":   "1185274747836174377",
	}
}
