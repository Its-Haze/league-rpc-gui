export interface ConnectionLightProps {
  label: string;
  connected: boolean;
}

// A small status dot: connected glows in the accent-adjacent ok color, idle
export function ConnectionLight({ label, connected }: ConnectionLightProps) {
  return (
    <span className="flex items-center gap-1.5 text-sm" title={`${label}: ${connected ? "connected" : "not connected"}`}>
      <span
        className={
          connected
            ? "bg-ok inline-block size-2 rounded-full"
            : "border-border inline-block size-2 rounded-full border"
        }
        style={connected ? { boxShadow: "0 0 6px var(--color-ok)" } : undefined}
        aria-hidden
      />
      <span className={connected ? "text-text" : "text-muted"}>{label}</span>
    </span>
  );
}
