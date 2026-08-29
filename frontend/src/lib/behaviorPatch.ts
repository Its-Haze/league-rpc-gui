import type { Config } from "../../bindings/github.com/its-haze/league-rpc/internal/config/models";

// Pure patch builders for the Behavior screen's toggles, each returning the
// Partial<Config> useSettings().applyPatch expects.

export function withLaunchAtStartup(cfg: Config, enabled: boolean): Partial<Config> {
  return { behavior: { ...cfg.behavior, launch_at_startup: enabled } };
}

/** What the window's close button does. Mirrors config.Close* in Go. */
export type CloseAction = "ask" | "tray" | "quit";

export function withCloseAction(cfg: Config, action: CloseAction): Partial<Config> {
  return { behavior: { ...cfg.behavior, close_action: action } };
}
