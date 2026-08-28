import { useSettings } from "../../hooks/useSettings";
import { useStatus } from "../../hooks/useStatus";
import { OnboardingWalkthrough } from "./home/OnboardingWalkthrough";
import { PresencePreview } from "./home/PresencePreview";
import { StatusDashboard } from "./home/StatusDashboard";

// The Home dashboard: live status, the last-sent presence preview, and the
// first-run walkthrough while onboarding_complete is still false.
export function HomeScreen() {
  const { cfg, applyPatch } = useSettings();
  const status = useStatus();

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-xl font-semibold">Home</h1>

      {cfg && !cfg.onboarding_complete && (
        <OnboardingWalkthrough
          cfg={cfg}
          status={status}
          applyPatch={applyPatch}
          onDismiss={() => void applyPatch({ onboarding_complete: true })}
        />
      )}

      <StatusDashboard status={status} />
      <PresencePreview status={status} />
    </div>
  );
}
