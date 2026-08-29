// Config.theme is a free-form string from the backend; narrow it defensively.
export const THEME_SETTINGS = ["system", "light", "dark"] as const;

export type ThemeSetting = (typeof THEME_SETTINGS)[number];
export type ResolvedTheme = "light" | "dark";

export function isThemeSetting(value: string): value is ThemeSetting {
  return (THEME_SETTINGS as readonly string[]).includes(value);
}

// Select-ready options for the theme picker: capitalized labels off the
// same source list, so a picker never drifts from what's actually valid.
export const THEME_OPTIONS = THEME_SETTINGS.map((value) => ({
  value,
  label: value[0].toUpperCase() + value.slice(1),
}));

// resolveTheme is the pure logic the applied-theme hook and its tests share:
// "system" tracks the OS, an explicit setting wins outright.
export function resolveTheme(setting: string, prefersDark: boolean): ResolvedTheme {
  if (setting === "light") return "light";
  if (setting === "dark") return "dark";
  return prefersDark ? "dark" : "light";
}
