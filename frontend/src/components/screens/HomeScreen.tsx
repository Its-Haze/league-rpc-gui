import { useSettings } from "../../hooks/useSettings";
import { useStatus } from "../../hooks/useStatus";
import { OnboardingWalkthrough } from "./home/OnboardingWalkthrough";
import { PresencePreview } from "./home/PresencePreview";
import { Welcome } from "./home/Welcome";

// The Home dashboard: a plain-language introduction, the first-run walkthrough,
// and the last-sent presence preview. No connection wiring here.
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

      <Welcome />
      <PresencePreview status={status} />
    </div>
  );
}
