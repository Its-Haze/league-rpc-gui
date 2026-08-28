// The presence contexts internal/presence/template knows about, display order.
export const PRESENCE_CONTEXTS = ["in-client", "champ-select", "in-game", "spectating"] as const;

export type PresenceContext = (typeof PRESENCE_CONTEXTS)[number];

export const PRESENCE_CONTEXT_LABELS: Record<PresenceContext, string> = {
  "in-client": "In client",
  "champ-select": "Champ select",
  "in-game": "In game",
  spectating: "Spectating",
};

export function isPresenceContext(value: string): value is PresenceContext {
  return (PRESENCE_CONTEXTS as readonly string[]).includes(value);
}
