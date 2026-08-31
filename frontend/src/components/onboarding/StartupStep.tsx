import { Monitor, Moon, Palette, Power, Sun } from "lucide-react";
import * as LabelPrimitive from "@radix-ui/react-label";
import type { Config } from "../../../bindings/github.com/its-haze/league-rpc/internal/config/models";
import { withLaunchAtStartup } from "../../lib/behaviorPatch";
import { THEME_SETTINGS, type ThemeSetting } from "../../lib/theme";
import { Toggle } from "../ui";

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
    <div className="flex flex-col gap-5 text-center">
      <h1 className="text-3xl font-semibold">Make it yours</h1>
      <p className="text-muted text-base">
        Two settings for the app itself, not the status it sends.
      </p>

      <section className="border-border bg-surface flex items-center justify-between gap-6 rounded-lg border p-5 text-left">
        <div className="flex items-start gap-3">
          <Palette className="text-accent mt-0.5 size-5 shrink-0" />
          <div className="flex flex-col gap-1">
            <h2 className="text-base font-semibold">Appearance</h2>
            <p className="text-muted text-sm">System follows whatever Windows is set to.</p>
          </div>
        </div>
        <div role="radiogroup" aria-label="Theme" className="border-border bg-surface-raised inline-flex shrink-0 items-center gap-1 rounded-lg border p-1">
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
      </section>

      <section className="border-accent/40 bg-accent/5 flex items-center justify-between gap-6 rounded-lg border p-5 text-left">
        <div className="flex items-start gap-3">
          <Power className="text-accent mt-0.5 size-5 shrink-0" />
          <div className="flex flex-col gap-1">
            <div className="flex flex-wrap items-center gap-2">
              <LabelPrimitive.Root htmlFor="onboarding-launch-at-startup" className="text-base font-semibold">
                Start with Windows
              </LabelPrimitive.Root>
              <span className="border-accent/40 text-accent rounded-full border px-2 py-0.5 text-xs font-medium">
                Recommended
              </span>
            </div>
            <p className="text-muted text-sm">
              Launches minimized to the tray, so your status is live before your first game.
            </p>
          </div>
        </div>
        <div className="shrink-0">
          <Toggle
            id="onboarding-launch-at-startup"
            checked={cfg.behavior.launch_at_startup}
            onCheckedChange={(v) => void applyPatch(withLaunchAtStartup(cfg, v))}
            label="Start with Windows"
          />
        </div>
      </section>
    </div>
  );
}
