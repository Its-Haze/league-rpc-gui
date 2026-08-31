import { CircleSlash, DownloadCloud, PanelTopClose, Power } from "lucide-react";
import { useEffect, useState } from "react";
import {
  GetStatus,
  SetPaused,
} from "../../../bindings/github.com/its-haze/league-rpc/cmd/league-rpc-gui/guiservice";
import { useCheckForUpdates } from "../../hooks/useCheckForUpdates";
import { useDefaultConfig } from "../../hooks/useDefaultConfig";
import { useSettings } from "../../hooks/useSettings";
import { useStatus } from "../../hooks/useStatus";
import { withCloseAction, withLaunchAtStartup, type CloseAction } from "../../lib/behaviorPatch";
import { Button, Select, SettingsCard, Toggle, type SelectOption } from "../ui";

const CLOSE_ACTIONS: SelectOption[] = [
  { value: "ask", label: "Ask me every time" },
  { value: "tray", label: "Hide to tray" },
  { value: "quit", label: "Quit League RPC" },
];

// The Behavior section: pausing presence, start-with-Windows and what closing
// the window does, and the update-check controls.
export function BehaviorScreen() {
  const { cfg, error, applyPatch } = useSettings();
  const defaults = useDefaultConfig();
  const { checking, result: checkResult, check: handleCheck } = useCheckForUpdates();
  const status = useStatus();
  const [paused, setPaused] = useState(false);

  // Local override wins until the daemon's own status:changed catches up, so
  // the toggle reflects the click immediately rather than the next broadcast.
  useEffect(() => {
    if (status) setPaused(status.paused);
  }, [status?.paused]);

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

  if (!cfg) {
    return <p className="text-muted text-sm">Loading settings…</p>;
  }

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-xl font-semibold">Behavior</h1>
      {error && <p className="text-danger text-sm">{error}</p>}

      <SettingsCard
        icon={CircleSlash}
        title="Pause presence"
        description="Stops updating your Discord status and clears it right away. Turn it back off whenever you like, and a restart unpauses it too."
        action={
          <Toggle
            id="pause-presence"
            checked={paused}
            onCheckedChange={togglePaused}
            label="Pause presence"
          />
        }
      />

      <SettingsCard
        icon={Power}
        title="Start with Windows"
        description="Launches minimized to the tray when you sign in, so your status is live before your first game."
        badge="Recommended"
        highlighted
        onReset={defaults ? () => void applyPatch(withLaunchAtStartup(cfg, defaults.behavior.launch_at_startup)) : undefined}
        isDefault={!defaults || cfg.behavior.launch_at_startup === defaults.behavior.launch_at_startup}
        action={
          <Toggle
            id="launch-at-startup"
            checked={cfg.behavior.launch_at_startup}
            onCheckedChange={(v) => void applyPatch(withLaunchAtStartup(cfg, v))}
            label="Start with Windows"
          />
        }
      />

      <SettingsCard
        icon={PanelTopClose}
        title="When I close the window"
        description="Hiding keeps your presence running in the background. Reopen it from the tray icon."
        onReset={defaults ? () => void applyPatch(withCloseAction(cfg, defaults.behavior.close_action as CloseAction)) : undefined}
        isDefault={!defaults || cfg.behavior.close_action === defaults.behavior.close_action}
        action={
          <Select
            aria-label="When I close the window"
            value={cfg.behavior.close_action}
            onValueChange={(v) => void applyPatch(withCloseAction(cfg, v as CloseAction))}
            options={CLOSE_ACTIONS}
          />
        }
      />

      <SettingsCard
        icon={DownloadCloud}
        title="Updates"
        description="You're on the stable channel. New versions install from inside the app."
        action={
          <>
            {checkResult && <span className="text-muted text-sm">{checkResult}</span>}
            <Button variant="secondary" onClick={handleCheck} disabled={checking}>
              {checking ? "Checking…" : "Check for updates"}
            </Button>
          </>
        }
      />
    </div>
  );
}
