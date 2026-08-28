import { describe, expect, it } from "vitest";
import { clearOverrideField, isOverridden } from "./displayOverride";

describe("clearOverrideField", () => {
  it("returns undefined for an already-undefined override", () => {
    expect(clearOverrideField(undefined, "show_rank")).toBeUndefined();
  });

  it("drops just the given field", () => {
    expect(clearOverrideField({ show_rank: false, show_stats: false }, "show_rank")).toEqual({
      show_stats: false,
    });
  });

  it("returns undefined once no fields remain overridden", () => {
    expect(clearOverrideField({ show_rank: false }, "show_rank")).toBeUndefined();
  });
});

describe("isOverridden", () => {
  it("is false when unset or undefined", () => {
    expect(isOverridden(undefined, "show_rank")).toBe(false);
    expect(isOverridden({}, "show_rank")).toBe(false);
  });

  it("is true once explicitly set, even to false", () => {
    expect(isOverridden({ show_rank: false }, "show_rank")).toBe(true);
  });
});
