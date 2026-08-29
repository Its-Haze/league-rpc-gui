import { useStatus } from "../../hooks/useStatus";
import { PresencePreview } from "./home/PresencePreview";

// The Home dashboard: the last-sent presence preview. No connection wiring
// here; the first-run flow (App Shell level) handles introducing the app.
export function HomeScreen() {
  const status = useStatus();

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-xl font-semibold">Home</h1>
      <PresencePreview status={status} />
    </div>
  );
}
