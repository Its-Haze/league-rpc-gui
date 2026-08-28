import { useReducer } from "react";
import type { Config } from "../../../../bindings/github.com/its-haze/league-rpc/internal/config/models";
import type { StatusSnapshot } from "../../../../bindings/github.com/its-haze/league-rpc/internal/app/models";
import { withShowRank, withShowStats } from "../../../lib/displayPatch";
import { connectStepReady, initialOnboarding, isFirstStep, isLastStep, onboardingReducer } from "../../../lib/onboarding";
import { Button, Toggle } from "../../ui";

export interface OnboardingWalkthroughProps {
  cfg: Config;
  status: StatusSnapshot | null;
  applyPatch: (patch: Partial<Config>) => Promise<void>;
  onDismiss: () => void;
}

// The first-run walkthrough: rendered inline on Home, not as a modal.
// Dismissible at any step; a Help entry can reopen it by clearing the flag.
export function OnboardingWalkthrough({ cfg, status, applyPatch, onDismiss }: OnboardingWalkthroughProps) {
  const [state, dispatch] = useReducer(onboardingReducer, initialOnboarding);
  const ready = connectStepReady(status?.discord_connected ?? false, status?.league_process ?? false);

  return (
    <section className="border-accent bg-surface-raised flex flex-col gap-4 rounded-lg border p-6">
      <div className="flex items-center justify-between">
        <h2 className="text-sm font-semibold">Welcome to League RPC</h2>
        <button onClick={onDismiss} className="text-muted text-xs underline">
          Skip
        </button>
      </div>

      {state.step === "welcome" && (
        <p className="text-sm">
          League RPC shows what you're doing in League of Legends as a Discord Rich Presence status:
          your queue, rank, and in-game stats, updated live.
        </p>
      )}

      {state.step === "connect" && (
        <div className="flex flex-col gap-2 text-sm">
          <p>Checking that League RPC can reach Discord and the League client:</p>
          <CheckItem label="Discord" done={status?.discord_connected ?? false} />
          <CheckItem label="League" done={status?.league_process ?? false} />
          {!ready && <p className="text-muted text-xs">You can continue setup while these connect.</p>}
        </div>
      )}

      {state.step === "display" && (
        <div className="flex flex-col gap-2 text-sm">
          <p>Pick what shows up in your presence:</p>
          <label className="flex items-center justify-between gap-4">
            Show rank
            <Toggle
              checked={cfg.display.default.show_rank}
              onCheckedChange={(v) => void applyPatch(withShowRank(cfg, v))}
              label="Show rank"
            />
          </label>
          <label className="flex items-center justify-between gap-4">
            Show stats
            <Toggle
              checked={cfg.display.default.show_stats}
              onCheckedChange={(v) => void applyPatch(withShowStats(cfg, v))}
              label="Show stats"
            />
          </label>
        </div>
      )}

      {state.step === "startup" && (
        <div className="flex flex-col gap-2 text-sm">
          <p>Start League RPC automatically when Windows starts?</p>
          <label className="flex items-center justify-between gap-4">
            Start with Windows
            <Toggle
              checked={cfg.behavior.launch_at_startup}
              onCheckedChange={(v) => void applyPatch({ behavior: { ...cfg.behavior, launch_at_startup: v } })}
              label="Start with Windows"
            />
          </label>
        </div>
      )}

      <div className="flex items-center justify-between">
        <Button variant="ghost" onClick={() => dispatch({ type: "back" })} disabled={isFirstStep(state.step)}>
          Back
        </Button>
        {isLastStep(state.step) ? (
          <Button variant="primary" onClick={onDismiss}>
            Finish
          </Button>
        ) : (
          <Button variant="primary" onClick={() => dispatch({ type: "next" })}>
            Next
          </Button>
        )}
      </div>
    </section>
  );
}

function CheckItem({ label, done }: { label: string; done: boolean }) {
  return (
    <div className="flex items-center gap-2">
      <span className={done ? "text-ok" : "text-muted"}>{done ? "✓" : "…"}</span>
      <span>{label}</span>
    </div>
  );
}
