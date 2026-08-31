import { THEME_ICONS, THEME_SETTINGS, type ThemeSetting } from "../../lib/theme";

export interface ThemePickerProps {
  value: string;
  onChange: (theme: ThemeSetting) => void;
}

// Segmented icon picker for the theme. Shared by onboarding and Behavior so
// the two never drift; the sidebar keeps a plain Select, which fits its width.
export function ThemePicker({ value, onChange }: ThemePickerProps) {
  return (
    <div
      role="radiogroup"
      aria-label="Theme"
      className="border-border bg-surface-raised inline-flex shrink-0 items-center gap-1 rounded-lg border p-1"
    >
      {THEME_SETTINGS.map((setting) => {
        const Icon = THEME_ICONS[setting];
        const active = value === setting;
        return (
          <button
            key={setting}
            type="button"
            role="radio"
            aria-checked={active}
            onClick={() => onChange(setting)}
            className={
              "flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium capitalize transition-colors " +
              (active ? "bg-surface text-accent" : "text-muted hover:text-text")
            }
          >
            <Icon className="size-4" />
            {setting}
          </button>
        );
      })}
    </div>
  );
}
