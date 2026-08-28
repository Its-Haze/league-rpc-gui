import type { ModeOverride } from "../../bindings/github.com/its-haze/league-rpc/internal/config/models";

export type OverrideField = "show_rank" | "show_stats";

// Clears field back to "inherit the default". Returns undefined once the
// mode has no fields left overridden, so an empty entry doesn't linger.
export function clearOverrideField(
  override: ModeOverride | undefined,
  field: OverrideField,
): ModeOverride | undefined {
  if (!override) return undefined;
  const next: ModeOverride = { ...override };
  delete next[field];
  return Object.keys(next).length > 0 ? next : undefined;
}

// Whether field is explicitly set on override (as opposed to inheriting).
export function isOverridden(override: ModeOverride | undefined, field: OverrideField): boolean {
  return override?.[field] != null;
}
