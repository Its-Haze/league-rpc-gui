import type { StatusSnapshot } from "../../../../bindings/github.com/its-haze/league-rpc/internal/app/models";
import { useDiscordAppName } from "../../../hooks/useDiscordAppName";
import { useSettings } from "../../../hooks/useSettings";
import { DiscordPresenceCard } from "../../DiscordPresenceCard";

export interface PresencePreviewProps {
  status: StatusSnapshot | null;
}

// Mirrors what the Updater actually last sent to Discord, never a
// recomputation, so this can never disagree with what Discord shows.
export function PresencePreview({ status }: PresencePreviewProps) {
  const { cfg } = useSettings();
  const appName = useDiscordAppName(cfg?.discord_app_id ?? "");
  const presence = status?.presence;

  return (
    <section className="border-border bg-surface flex flex-col gap-3 rounded-lg border p-6">
      <h2 className="text-sm font-semibold">Discord presence</h2>

      {status?.presence_cleared || !presence ? (
        <p className="text-muted text-sm">Presence is currently cleared.</p>
      ) : (
        <DiscordPresenceCard
          details={presence.Details}
          state={presence.State}
          largeImage={presence.LargeImage}
          largeText={presence.LargeText}
          smallImage={presence.SmallImage}
          smallText={presence.SmallText}
          appName={appName}
          startUnixSeconds={presence.Start}
        />
      )}
    </section>
  );
}
