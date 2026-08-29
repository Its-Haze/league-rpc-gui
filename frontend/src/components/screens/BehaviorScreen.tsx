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
import { Field, Select, Toggle, type SelectOption } from "../ui";

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

      <section className="border-border bg-surface flex flex-col gap-1 rounded-lg border p-6">
        <Field
          id="pause-presence"
          label="Pause presence"
          hint="Clears your Discord status immediately, resumes on next launch"
        >
          <Toggle
            id="pause-presence"
            checked={paused}
            onCheckedChange={togglePaused}
            label="Pause presence"
          />
        </Field>
      </section>

      <section className="border-border bg-surface flex flex-col gap-1 rounded-lg border p-6">
        <Field
          id="launch-at-startup"
          label="Start with Windows"
          hint="Launches minimized to the tray"
          onReset={defaults ? () => void applyPatch(withLaunchAtStartup(cfg, defaults.behavior.launch_at_startup)) : undefined}
          isDefault={!defaults || cfg.behavior.launch_at_startup === defaults.behavior.launch_at_startup}
        >
          <Toggle
            id="launch-at-startup"
            checked={cfg.behavior.launch_at_startup}
            onCheckedChange={(v) => void applyPatch(withLaunchAtStartup(cfg, v))}
            label="Start with Windows"
          />
        </Field>
        <Field
          id="close-action"
          label="When I close the window"
          hint="Hiding keeps presence running. Reopen from the tray icon"
          onReset={defaults ? () => void applyPatch(withCloseAction(cfg, defaults.behavior.close_action as CloseAction)) : undefined}
          isDefault={!defaults || cfg.behavior.close_action === defaults.behavior.close_action}
        >
          <Select
            aria-label="When I close the window"
            value={cfg.behavior.close_action}
            onValueChange={(v) => void applyPatch(withCloseAction(cfg, v as CloseAction))}
            options={CLOSE_ACTIONS}
          />
        </Field>
      </section>

      <section className="border-border bg-surface flex flex-col gap-2 rounded-lg border p-6">
        <div className="flex items-center justify-between gap-4">
          <div className="flex flex-col gap-0.5">
            <span className="text-sm">Update channel</span>
            <span className="text-muted text-xs">Stable only for now</span>
          </div>
          <span className="text-muted text-sm">stable</span>
        </div>
        <div className="flex items-center gap-3">
          <button
            onClick={handleCheck}
            disabled={checking}
            className="border-border bg-surface-raised rounded-sm border px-3 py-1.5 text-sm disabled:opacity-60"
          >
            {checking ? "Checking…" : "Check for updates"}
          </button>
          {checkResult && <span className="text-muted text-sm">{checkResult}</span>}
        </div>
      </section>
    </div>
  );
}
