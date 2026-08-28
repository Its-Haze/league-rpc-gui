import { useCheckForUpdates } from "../../hooks/useCheckForUpdates";
import { useSettings } from "../../hooks/useSettings";
import { withIdleText, withLaunchAtStartup, withShowInClient } from "../../lib/behaviorPatch";
import { DebouncedTextField, Field, Toggle } from "../ui";

// The Behavior section: start-with-Windows, the tray explainer, in-client
// and idle presence, and the update-check controls.
export function BehaviorScreen() {
  const { cfg, error, applyPatch } = useSettings();
  const { checking, result: checkResult, check: handleCheck } = useCheckForUpdates();

  if (!cfg) {
    return <p className="text-muted text-sm">Loading settings…</p>;
  }

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-xl font-semibold">Behavior</h1>
      {error && <p className="text-danger text-sm">{error}</p>}

      <section className="border-border bg-surface flex flex-col gap-1 rounded-lg border p-6">
        <Field id="launch-at-startup" label="Start with Windows" hint="Launches hidden to the tray at logon">
          <Toggle
            id="launch-at-startup"
            checked={cfg.behavior.launch_at_startup}
            onCheckedChange={(v) => void applyPatch(withLaunchAtStartup(cfg, v))}
            label="Start with Windows"
          />
        </Field>
        <p className="text-muted text-xs">
          Closing the window hides League RPC to the system tray; presence keeps running. Right-click
          the tray icon and choose Quit to stop it fully.
        </p>
      </section>

      <section className="border-border bg-surface flex flex-col gap-1 rounded-lg border p-6">
        <Field id="show-in-client" label="Show presence while in client" hint="Idle in the League client, not in a game">
          <Toggle
            id="show-in-client"
            checked={cfg.presence.show_in_client}
            onCheckedChange={(v) => void applyPatch(withShowInClient(cfg, v))}
            label="Show presence while in client"
          />
        </Field>
        <Field id="idle-text" label="Idle status text" hint="Empty uses the built-in text">
          <DebouncedTextField
            id="idle-text"
            value={cfg.presence.idle}
            onCommit={(v) => void applyPatch(withIdleText(cfg, v))}
            className="border-border bg-surface-raised text-text rounded-sm border px-3 py-1.5 text-sm"
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
