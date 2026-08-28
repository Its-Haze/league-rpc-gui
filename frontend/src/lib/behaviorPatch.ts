import type { Config } from "../../bindings/github.com/its-haze/league-rpc/internal/config/models";

// Pure patch builders for the Behavior screen's toggles, each returning the
// Partial<Config> useSettings().applyPatch expects.

export function withLaunchAtStartup(cfg: Config, enabled: boolean): Partial<Config> {
  return { behavior: { ...cfg.behavior, launch_at_startup: enabled } };
}

export function withShowInClient(cfg: Config, enabled: boolean): Partial<Config> {
  return { presence: { ...cfg.presence, show_in_client: enabled } };
}

export function withIdleText(cfg: Config, idle: string): Partial<Config> {
  return { presence: { ...cfg.presence, idle } };
}
