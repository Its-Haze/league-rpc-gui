package config

import "encoding/json"

// legacyConfig is the flat, versionless schema (v1). It exists only so a file
// written by an older build can be read and mapped onto the current tree.
type legacyConfig struct {
	DiscordAppID         string `json:"discord_app_id"`
	ShowStats            bool   `json:"show_stats"`
	ShowRank             bool   `json:"show_rank"`
	ShowEmojis           bool   `json:"show_emojis"`
	ShowInClient         bool   `json:"show_in_client"`
	UpdateInterval       int    `json:"update_interval"`
	StatsPollingInterval int    `json:"stats_polling_interval"`
	DebugMode            bool   `json:"debug_mode"`
}

// schemaVersionOf reports the schema_version recorded in raw. A missing or
// unparseable field reads as 0, which routes the file through migration.
func schemaVersionOf(raw []byte) int {
	var probe struct {
		SchemaVersion int `json:"schema_version"`
	}
	_ = json.Unmarshal(raw, &probe)
	return probe.SchemaVersion
}

// migrateFromV1 parses raw as the flat schema and folds it into a default
// tree, so any field the old file omitted keeps its default.
func migrateFromV1(raw []byte) (*Config, error) {
	var l legacyConfig
	if err := json.Unmarshal(raw, &l); err != nil {
		return nil, err
	}

	c := DefaultConfig()
	c.DiscordAppID = l.DiscordAppID
	c.Display.Default.ShowRank = l.ShowRank
	c.Display.Default.ShowStats = l.ShowStats
	c.Presence.ShowEmojis = l.ShowEmojis
	c.Presence.ShowInClient = l.ShowInClient
	c.Advanced.UpdateInterval = l.UpdateInterval
	c.Advanced.StatsPollingInterval = l.StatsPollingInterval
	c.Advanced.DebugMode = l.DebugMode
	c.SchemaVersion = CurrentSchemaVersion
	// An upgrading install already went through first-run setup before
	// onboarding existed; don't show it the walkthrough again.
	c.OnboardingComplete = true
	return c, nil
}

// migrateV2ToV3 parses raw as a v2 tree, same shape as Config but possibly
// missing OnboardingComplete, and backfills it.
func migrateV2ToV3(raw []byte) (*Config, error) {
	var c Config
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, err
	}
	c.OnboardingComplete = true
	c.SchemaVersion = CurrentSchemaVersion
	return &c, nil
}
