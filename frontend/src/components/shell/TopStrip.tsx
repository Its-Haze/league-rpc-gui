import { useEffect, useState } from "react";
import { Events } from "@wailsio/runtime";
import {
  GetStatus,
  SetPaused,
} from "../../../bindings/github.com/its-haze/league-rpc/cmd/league-rpc-gui/guiservice";
import type { StatusSnapshot } from "../../../bindings/github.com/its-haze/league-rpc/internal/app/models";
import { Select, Toggle } from "../ui";
import { THEME_SETTINGS, type ThemeSetting } from "../../lib/theme";
import { ConnectionLight } from "./ConnectionLight";

const STATUS_CHANGED_EVENT = "status:changed";

const THEME_OPTIONS = THEME_SETTINGS.map((value) => ({
  value,
  label: value[0].toUpperCase() + value.slice(1),
}));

export interface TopStripProps {
  theme: string;
  onThemeChange: (theme: ThemeSetting) => void;
  /** True until the initial settings load resolves, so the picker can't fire on a config that isn't there yet. */
  themeDisabled?: boolean;
}

// Persistent strip shown above every screen: the three connection lights
// driven by status:changed, the Pause toggle, and the theme picker.
export function TopStrip({ theme, onThemeChange, themeDisabled }: TopStripProps) {
  const [status, setStatus] = useState<StatusSnapshot | null>(null);
  const [paused, setPaused] = useState(false);

  useEffect(() => {
    GetStatus().then((s) => {
      setStatus(s);
      setPaused(s.paused);
    }).catch(() => {});

    const off = Events.On(STATUS_CHANGED_EVENT, (ev: { data: StatusSnapshot }) => {
      setStatus(ev.data);
      setPaused(ev.data.paused);
    });
    return () => off();
  }, []);

  async function togglePaused(next: boolean) {
    setPaused(next);
    try {
      await SetPaused(next);
    } catch {
      // Re-sync from the daemon rather than assuming !next: a status:changed
      // broadcast may have already landed while this call was in flight.
      GetStatus()
        .then((s) => setPaused(s.paused))
        .catch(() => setPaused(!next));
    }
  }

  return (
    <div className="border-border bg-surface flex items-center justify-between gap-6 border-b px-6 py-3">
      <div className="flex items-center gap-5">
        <ConnectionLight label="League" connected={status?.league_process ?? false} />
        <ConnectionLight label="League Client" connected={status?.lcu_connected ?? false} />
        <ConnectionLight label="Discord" connected={status?.discord_connected ?? false} />
      </div>

      <div className="flex items-center gap-5">
        <label className="flex items-center gap-2 text-sm">
          <span className="text-muted">Pause presence</span>
          <Toggle checked={paused} onCheckedChange={togglePaused} label="Pause presence" />
        </label>

        <Select
          value={theme}
          onValueChange={(v) => onThemeChange(v as ThemeSetting)}
          options={THEME_OPTIONS}
          disabled={themeDisabled}
          aria-label="Theme"
        />
      </div>
    </div>
  );
}
