import { Monitor, Moon, Sun } from "lucide-react";
import type { Config } from "../../../bindings/github.com/its-haze/league-rpc/internal/config/models";
import { withLaunchAtStartup } from "../../lib/behaviorPatch";
import { THEME_SETTINGS, type ThemeSetting } from "../../lib/theme";
import { Field, Toggle } from "../ui";

export interface StartupStepProps {
  cfg: Config;
  applyPatch: (patch: Partial<Config>) => Promise<void>;
}

const THEME_ICONS: Record<ThemeSetting, typeof Monitor> = {
  system: Monitor,
  light: Sun,
  dark: Moon,
};

// Onboarding's third screen: settings about the app itself, not the presence it sends.
export function StartupStep({ cfg, applyPatch }: StartupStepProps) {
  return (
    <div className="flex flex-col gap-4 text-center">
      <h1 className="text-2xl font-semibold">Make it yours</h1>
      <p className="text-muted text-sm">
        A couple of settings for the app itself, not your presence.
      </p>

      <section className="border-border bg-surface flex flex-col gap-3 rounded-lg border p-4">
        <div className="flex items-center justify-between gap-4">
          <div className="flex flex-col gap-0.5 text-left">
            <span className="text-sm">Theme</span>
          </div>
          <div role="radiogroup" aria-label="Theme" className="border-border bg-surface-raised inline-flex items-center gap-1 rounded-lg border p-1">
            {THEME_SETTINGS.map((setting) => {
              const Icon = THEME_ICONS[setting];
              const active = cfg.theme === setting;
              return (
                <button
                  key={setting}
                  role="radio"
                  aria-checked={active}
                  onClick={() => void applyPatch({ theme: setting })}
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
        </div>
      </section>

      <section className="border-border bg-surface flex flex-col gap-1 rounded-lg border p-4 text-left">
        <Field id="onboarding-launch-at-startup" label="Start with Windows" hint="Launches minimized to the tray">
          <Toggle
            id="onboarding-launch-at-startup"
            checked={cfg.behavior.launch_at_startup}
            onCheckedChange={(v) => void applyPatch(withLaunchAtStartup(cfg, v))}
            label="Start with Windows"
          />
        </Field>
        <p className="text-muted text-xs">
          We recommend keeping this on, so your Discord status is always ready to show off your
          games, without you having to remember to start League RPC yourself.
        </p>
      </section>
    </div>
  );
}
