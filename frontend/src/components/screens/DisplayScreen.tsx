import { useEffect, useState } from "react";
import { GetGameModes } from "../../../bindings/github.com/its-haze/league-rpc/cmd/league-rpc-gui/guiservice";
import type { ModeOverride, TemplatePair } from "../../../bindings/github.com/its-haze/league-rpc/internal/config/models";
import { useSettings } from "../../hooks/useSettings";
import { withShowEmojis, withShowRank, withShowStats } from "../../lib/displayPatch";
import { PRESENCE_CONTEXTS } from "../../lib/presenceContexts";
import { Field, Toggle } from "../ui";
import { ModeOverrideRow } from "./display/ModeOverrideRow";
import { TemplateEditor } from "./display/TemplateEditor";

// The Display section: global toggles, per-GameMode overrides, and the
// per-context presence template editors.
export function DisplayScreen() {
  const { cfg, error, applyPatch } = useSettings();
  const [modes, setModes] = useState<string[]>([]);
  const [modesOpen, setModesOpen] = useState(false);

  useEffect(() => {
    GetGameModes()
      .then((m) => setModes(m ?? []))
      .catch(() => {});
  }, []);

  if (!cfg) {
    return <p className="text-muted text-sm">Loading settings…</p>;
  }

  function setModeOverride(mode: string, next: ModeOverride | undefined) {
    const modesMap = { ...cfg!.display.modes };
    if (next) {
      modesMap[mode] = next;
    } else {
      delete modesMap[mode];
    }
    void applyPatch({ display: { ...cfg!.display, modes: modesMap } });
  }

  function setTemplate(ctx: string, next: TemplatePair) {
    void applyPatch({
      presence: { ...cfg!.presence, templates: { ...cfg!.presence.templates, [ctx]: next } },
    });
  }

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-xl font-semibold">Display</h1>
      {error && <p className="text-danger text-sm">{error}</p>}

      <section className="border-border bg-surface flex flex-col gap-1 rounded-lg border p-6">
        <Field id="show-rank" label="Show rank" hint="Rank emblem and LP">
          <Toggle
            id="show-rank"
            checked={cfg.display.default.show_rank}
            onCheckedChange={(v) => void applyPatch(withShowRank(cfg, v))}
            label="Show rank"
          />
        </Field>
        <Field id="show-stats" label="Show stats" hint="KDA and creep score">
          <Toggle
            id="show-stats"
            checked={cfg.display.default.show_stats}
            onCheckedChange={(v) => void applyPatch(withShowStats(cfg, v))}
            label="Show stats"
          />
        </Field>
        <Field id="show-emojis" label="Show status emojis" hint="Online/away indicator">
          <Toggle
            id="show-emojis"
            checked={cfg.presence.show_emojis}
            onCheckedChange={(v) => void applyPatch(withShowEmojis(cfg, v))}
            label="Show status emojis"
          />
        </Field>
      </section>

      <section className="border-border bg-surface flex flex-col gap-2 rounded-lg border p-6">
        <button
          onClick={() => setModesOpen((v) => !v)}
          className="flex items-center justify-between text-left text-sm font-semibold"
        >
          <span>Per-mode overrides</span>
          <span className="text-muted">{modesOpen ? "▾" : "▸"}</span>
        </button>
        {modesOpen && (
          <div className="flex flex-col">
            {modes.length === 0 ? (
              <p className="text-muted text-sm">Loading game modes…</p>
            ) : (
              modes.map((mode) => (
                <ModeOverrideRow
                  key={mode}
                  mode={mode}
                  def={cfg.display.default}
                  override={cfg.display.modes?.[mode]}
                  onChange={(next) => setModeOverride(mode, next)}
                />
              ))
            )}
          </div>
        )}
      </section>

      <section className="flex flex-col gap-4">
        <h2 className="text-sm font-semibold">Presence text</h2>
        {PRESENCE_CONTEXTS.map((ctx) => {
          const pair = cfg.presence.templates?.[ctx] ?? { details: "", state: "" };
          return (
            <TemplateEditor
              key={ctx}
              ctx={ctx}
              value={pair}
              onChange={(next) => setTemplate(ctx, next)}
              showRank={cfg.display.default.show_rank}
              showStats={cfg.display.default.show_stats}
            />
          );
        })}
      </section>
    </div>
  );
}
