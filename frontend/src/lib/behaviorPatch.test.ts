import { describe, expect, it } from "vitest";
import { DefaultConfig } from "./testFixtures";
import { withIdleText, withLaunchAtStartup, withShowInClient } from "./behaviorPatch";

describe("withLaunchAtStartup", () => {
  it("sets the flag and keeps sibling behavior fields", () => {
    const cfg = DefaultConfig();
    const patch = withLaunchAtStartup(cfg, true);
    expect(patch.behavior).toEqual({ ...cfg.behavior, launch_at_startup: true });
  });
});

describe("withShowInClient", () => {
  it("sets show_in_client and keeps sibling presence fields", () => {
    const cfg = DefaultConfig();
    const patch = withShowInClient(cfg, false);
    expect(patch.presence).toEqual({ ...cfg.presence, show_in_client: false });
  });
});

describe("withIdleText", () => {
  it("sets the idle override text", () => {
    const cfg = DefaultConfig();
    const patch = withIdleText(cfg, "Taking a break");
    expect(patch.presence).toEqual({ ...cfg.presence, idle: "Taking a break" });
  });
});
