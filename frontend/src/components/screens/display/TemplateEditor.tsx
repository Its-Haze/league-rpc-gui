import { useEffect, useRef, useState } from "react";
import {
  GetDisplayPreview,
  GetTemplateTokens,
} from "../../../../bindings/github.com/its-haze/league-rpc/cmd/league-rpc-gui/guiservice";
import type { TemplatePair } from "../../../../bindings/github.com/its-haze/league-rpc/internal/config/models";
import { useDebouncedValue } from "../../../hooks/useDebouncedValue";
import { usePreviewAssets } from "../../../hooks/usePreviewAssets";
import { PRESENCE_CONTEXT_LABELS, type PresenceContext } from "../../../lib/presenceContexts";
import { DiscordPresenceCard } from "../../DiscordPresenceCard";
import { Field } from "../../ui";

// How long to wait after the last keystroke before persisting the template
// and re-running the preview, so typing doesn't save/round-trip every key.
const COMMIT_DELAY_MS = 400;

export interface TemplateEditorProps {
  ctx: PresenceContext;
  value: TemplatePair;
  onChange: (next: TemplatePair) => void;
  /** Current Display.Default toggles, so the preview honors them the same
   * way a real send would (stats line, rank emblem vs. the League logo). */
  showRank: boolean;
  showStats: boolean;
}

// One presence context's editor: details/state text fields, a live preview
// rendered through the real template engine, and a token reference.
export function TemplateEditor({ ctx, value, onChange, showRank, showStats }: TemplateEditorProps) {
  const [tokens, setTokens] = useState<string[]>([]);
  const [preview, setPreview] = useState<{ details: string; state: string; warnings: string[] }>({
    details: "",
    state: "",
    warnings: [],
  });
  const assets = usePreviewAssets();

  // Local draft so typing feels instant; onChange (which persists to the
  // daemon) only fires once the draft settles, via the debounce below.
  const [draft, setDraft] = useState(value);
  useEffect(() => setDraft(value), [value]);
  const debouncedDraft = useDebouncedValue(draft, COMMIT_DELAY_MS);
  const mounted = useRef(false);

  useEffect(() => {
    GetTemplateTokens(ctx)
      .then((t) => setTokens(t ?? []))
      .catch(() => {});
  }, [ctx]);

  useEffect(() => {
    if (!mounted.current) {
      mounted.current = true;
    } else {
      onChange(debouncedDraft);
    }

    let cancelled = false;
    GetDisplayPreview(ctx, debouncedDraft, showStats)
      .then((p) => {
        if (!cancelled) setPreview({ details: p.details, state: p.state, warnings: p.warnings ?? [] });
      })
      .catch(() => {});
    return () => {
      cancelled = true;
    };
    // onChange is expected to be stable enough per render; only ctx, the
  }, [ctx, debouncedDraft, showStats]);

  // Sample images only make sense for the in-game context: it's the only
  // one whose sample data includes a champion/skin to illustrate.
  const showImages = ctx === "in-game" && assets;
  const largeImage = showImages ? assets.champion_skin_url : undefined;
  const smallImage = showImages ? (showRank ? assets.rank_emblem_url : assets.league_logo_url) : undefined;

  return (
    <div className="border-border flex flex-col gap-3 rounded-md border p-4">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-semibold">{PRESENCE_CONTEXT_LABELS[ctx]}</h3>
        <span className="text-muted text-xs">
          Tokens: {tokens.length > 0 ? tokens.map((t) => `{${t}}`).join(" ") : "none"}
        </span>
      </div>

      <Field id={`${ctx}-details`} label="Details line">
        <input
          id={`${ctx}-details`}
          value={draft.details}
          onChange={(e) => setDraft({ ...draft, details: e.target.value })}
          placeholder="blank uses the built-in default"
          className="border-border bg-surface-raised text-text w-64 rounded-sm border px-3 py-1.5 text-sm"
        />
      </Field>
      <Field id={`${ctx}-state`} label="State line">
        <input
          id={`${ctx}-state`}
          value={draft.state}
          onChange={(e) => setDraft({ ...draft, state: e.target.value })}
          placeholder="blank uses the built-in default"
          className="border-border bg-surface-raised text-text w-64 rounded-sm border px-3 py-1.5 text-sm"
        />
      </Field>

      <div className="flex flex-col gap-1">
        <div className="text-muted text-xs">Preview</div>
        <DiscordPresenceCard details={preview.details} state={preview.state} largeImage={largeImage} smallImage={smallImage} />
      </div>

      {preview.warnings.length > 0 && (
        <ul className="text-danger text-xs">
          {preview.warnings.map((w) => (
            <li key={w}>{w}</li>
          ))}
        </ul>
      )}
    </div>
  );
}
