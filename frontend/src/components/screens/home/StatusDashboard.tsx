import type { StatusSnapshot } from "../../../../bindings/github.com/its-haze/league-rpc/internal/app/models";
import { ConnectionLight } from "../../shell/ConnectionLight";

const PHASE_LABELS: Record<string, string> = {
  None: "Idle in client",
  Lobby: "In a lobby",
  Matchmaking: "In queue",
  ReadyCheck: "Ready check",
  ChampSelect: "Champion select",
  GameStart: "Loading into game",
  InProgress: "In game",
  Watching: "Spectating",
  Reconnect: "Reconnecting",
  WaitingForStats: "Waiting for stats",
  PreEndOfGame: "Game just ended",
  EndOfGame: "Post-game",
  TerminatedInError: "Game ended in error",
};

export interface StatusDashboardProps {
  status: StatusSnapshot | null;
}

// The three connection states with plain-language explanations, plus the
// current game flow phase.
export function StatusDashboard({ status }: StatusDashboardProps) {
  return (
    <section className="border-border bg-surface flex flex-col gap-4 rounded-lg border p-6">
      <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
        <StatusRow
          label="League"
          connected={status?.league_process ?? false}
          okText="The League client is running."
          badText="League isn't running. Presence stays cleared until it starts."
        />
        <StatusRow
          label="League Client"
          connected={status?.lcu_connected ?? false}
          okText="Connected to the League Client API."
          badText="Waiting for the League Client API."
        />
        <StatusRow
          label="Discord"
          connected={status?.discord_connected ?? false}
          okText="Connected to Discord."
          badText="Discord isn't reachable. Presence won't show until it's running."
        />
      </div>
      <div className="text-muted text-sm">
        Phase: {PHASE_LABELS[status?.gameflow_phase ?? ""] ?? status?.gameflow_phase ?? "Unknown"}
      </div>
    </section>
  );
}

function StatusRow({
  label,
  connected,
  okText,
  badText,
}: {
  label: string;
  connected: boolean;
  okText: string;
  badText: string;
}) {
  return (
    <div className="flex flex-col gap-1">
      <ConnectionLight label={label} connected={connected} />
      <p className="text-muted text-xs">{connected ? okText : badText}</p>
    </div>
  );
}
