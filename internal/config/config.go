package config

import (
	"errors"
	"fmt"

	"github.com/its-haze/league-rpc/internal/presence/template"
)

// CurrentSchemaVersion is the version stamped on every config the app writes.
// A file with a lower (or missing) version is migrated on Load.
const CurrentSchemaVersion = 2

// Config is the versioned settings tree. It is persisted as JSON and edited
// through the GUI; nothing reads settings from CLI flags.
type Config struct {
	SchemaVersion int    `json:"schema_version"`
	DiscordAppID  string `json:"discord_app_id"`
	Theme         string `json:"theme"` // system | light | dark

	Display  DisplayConfig  `json:"display"`
	Presence PresenceConfig `json:"presence"`
	Behavior BehaviorConfig `json:"behavior"`
	Advanced AdvancedConfig `json:"advanced"`
}

// DisplayConfig holds what presence shows, with per-GameMode overrides layered
// on top of the global defaults.
type DisplayConfig struct {
	Default DisplayDefaults         `json:"default"`
	Modes   map[string]ModeOverride `json:"modes"`
}

// DisplayDefaults is the global on/off state for the two per-mode-overridable
// toggles.
type DisplayDefaults struct {
	ShowRank  bool `json:"show_rank"`  // rank emblem and LP
	ShowStats bool `json:"show_stats"` // KDA and creep score
}

// ModeOverride is a per-GameMode override. A nil field means "inherit the
// default"; a non-nil field means the user set it explicitly, false included.
type ModeOverride struct {
	ShowRank  *bool `json:"show_rank,omitempty"`
	ShowStats *bool `json:"show_stats,omitempty"`
}

// Resolve returns the effective toggles for mode, falling back to the default
// for any field the mode does not override.
func (d DisplayConfig) Resolve(mode string) DisplayDefaults {
	out := d.Default
	if ov, ok := d.Modes[mode]; ok {
		if ov.ShowRank != nil {
			out.ShowRank = *ov.ShowRank
		}
		if ov.ShowStats != nil {
			out.ShowStats = *ov.ShowStats
		}
	}
	return out
}

// PresenceConfig holds presence-wide text settings.
type PresenceConfig struct {
	ShowEmojis   bool                    `json:"show_emojis"`    // online/away emoji
	ShowInClient bool                    `json:"show_in_client"` // presence while idle in client
	Idle         string                  `json:"idle"`           // idle status override; empty uses the built-in
	Templates    map[string]TemplatePair `json:"templates"`      // per-context text, keyed by context
}

// TemplatePair is the editable text for one presence context: the two lines
// Discord renders.
type TemplatePair struct {
	Details string `json:"details"`
	State   string `json:"state"`
}

// BehaviorConfig holds launch-related settings.
type BehaviorConfig struct {
	LaunchAtStartup  bool   `json:"launch_at_startup"`  // start with Windows
	AutoLaunchLeague bool   `json:"auto_launch_league"` // launch League when the app starts
	LeaguePath       string `json:"league_path"`        // custom League install path; empty auto-detects
}

// AdvancedConfig holds tuning knobs and debug options.
type AdvancedConfig struct {
	UpdateInterval       int  `json:"update_interval"`        // RPC update throttle, ms
	StatsPollingInterval int  `json:"stats_polling_interval"` // in-game stats poll, ms
	DebugMode            bool `json:"debug_mode"`             // verbose logging
}

// DefaultConfig returns a fully populated tree at the current schema version.
func DefaultConfig() *Config {
	return &Config{
		SchemaVersion: CurrentSchemaVersion,
		DiscordAppID:  "1237146703111393281", // TEMP dev app ID; real one is 1194034071588851783
		Theme:         ThemeSystem,
		Display: DisplayConfig{
			Default: DisplayDefaults{ShowRank: true, ShowStats: true},
			Modes:   map[string]ModeOverride{},
		},
		Presence: PresenceConfig{
			ShowEmojis:   true,
			ShowInClient: true,
			Idle:         "",
			Templates:    defaultTemplates(),
		},
		Behavior: BehaviorConfig{
			LaunchAtStartup:  false,
			AutoLaunchLeague: false,
			LeaguePath:       "",
		},
		Advanced: AdvancedConfig{
			UpdateInterval:       1500,
			StatsPollingInterval: 3000,
			DebugMode:            false,
		},
	}
}

// defaultTemplates returns the built-in per-context presence templates. Every
// entry renders to the string the app produced before templates were editable.
func defaultTemplates() map[string]TemplatePair {
	m := make(map[string]TemplatePair, len(template.Contexts()))
	for _, ctx := range template.Contexts() {
		d, s := template.Default(ctx)
		m[string(ctx)] = TemplatePair{Details: d, State: s}
	}
	return m
}

// Bounds for the numeric settings. Shared by Validate and clamp so the
// reject path and the load-time repair path can't drift apart.
const (
	MinUpdateInterval       = 500
	MaxUpdateInterval       = 10000
	MinStatsPollingInterval = 1000
	MaxStatsPollingInterval = 30000
)

// Theme values.
const (
	ThemeSystem = "system"
	ThemeLight  = "light"
	ThemeDark   = "dark"
)

func validTheme(t string) bool {
	return t == ThemeSystem || t == ThemeLight || t == ThemeDark
}

// Validate reports every problem with c without changing it. Store.Apply and
// Save call this and surface the error; the GUI shows it next to the field.
func (c *Config) Validate() error {
	var errs []error

	if c.DiscordAppID == "" {
		errs = append(errs, errors.New("discord_app_id must not be empty"))
	}
	if !validTheme(c.Theme) {
		errs = append(errs, fmt.Errorf("theme must be one of %q, %q, %q", ThemeSystem, ThemeLight, ThemeDark))
	}
	if c.Advanced.UpdateInterval < MinUpdateInterval || c.Advanced.UpdateInterval > MaxUpdateInterval {
		errs = append(errs, fmt.Errorf("update_interval must be between %d and %d ms", MinUpdateInterval, MaxUpdateInterval))
	}
	if c.Advanced.StatsPollingInterval < MinStatsPollingInterval || c.Advanced.StatsPollingInterval > MaxStatsPollingInterval {
		errs = append(errs, fmt.Errorf("stats_polling_interval must be between %d and %d ms", MinStatsPollingInterval, MaxStatsPollingInterval))
	}

	return errors.Join(errs...)
}

// clamp forces c into valid bounds in place. Only the load path uses it, so a
// hand-broken or older config file on disk still boots with sane values.
func (c *Config) clamp() {
	def := DefaultConfig()

	if c.DiscordAppID == "" {
		c.DiscordAppID = def.DiscordAppID
	}
	if !validTheme(c.Theme) {
		c.Theme = def.Theme
	}
	if c.Advanced.UpdateInterval < MinUpdateInterval || c.Advanced.UpdateInterval > MaxUpdateInterval {
		c.Advanced.UpdateInterval = def.Advanced.UpdateInterval
	}
	if c.Advanced.StatsPollingInterval < MinStatsPollingInterval {
		c.Advanced.StatsPollingInterval = MinStatsPollingInterval
	}
	if c.Advanced.StatsPollingInterval > MaxStatsPollingInterval {
		c.Advanced.StatsPollingInterval = MaxStatsPollingInterval
	}
	if c.Display.Modes == nil {
		c.Display.Modes = map[string]ModeOverride{}
	}
	if c.Presence.Templates == nil {
		c.Presence.Templates = map[string]TemplatePair{}
	}
	// Backfill any context the file is missing so the GUI always has an entry
	// to show and edit; never touch one the user already set.
	for k, v := range defaultTemplates() {
		if _, ok := c.Presence.Templates[k]; !ok {
			c.Presence.Templates[k] = v
		}
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
