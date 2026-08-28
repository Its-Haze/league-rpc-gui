import type {
  DisplayDefaults,
  ModeOverride,
} from "../../../../bindings/github.com/its-haze/league-rpc/internal/config/models";
import { clearOverrideField, isOverridden, type OverrideField } from "../../../lib/displayOverride";
import { Select } from "../../ui";

const OPTIONS = [
  { value: "default", label: "Default" },
  { value: "on", label: "On" },
  { value: "off", label: "Off" },
];

function selectValue(override: ModeOverride | undefined, field: OverrideField): string {
  if (!isOverridden(override, field)) return "default";
  return override![field] ? "on" : "off";
}

export interface ModeOverrideRowProps {
  mode: string;
  def: DisplayDefaults;
  override: ModeOverride | undefined;
  onChange: (next: ModeOverride | undefined) => void;
}

// One GameMode row: a default/on/off control for each per-mode-overridable
// field. "Default" clears the field back to inheriting DisplayDefaults.
export function ModeOverrideRow({ mode, def, override, onChange }: ModeOverrideRowProps) {
  function setField(field: OverrideField, value: string) {
    if (value === "default") {
      onChange(clearOverrideField(override, field));
      return;
    }
    onChange({ ...override, [field]: value === "on" });
  }

  return (
    <div className="border-border flex items-center justify-between gap-4 border-t py-2 first:border-t-0">
      <span className="text-sm">{mode}</span>
      <div className="flex items-center gap-3">
        <label className="text-muted flex items-center gap-2 text-xs">
          Rank {!isOverridden(override, "show_rank") && `(${def.show_rank ? "on" : "off"})`}
          <Select
            value={selectValue(override, "show_rank")}
            onValueChange={(v) => setField("show_rank", v)}
            options={OPTIONS}
            aria-label={`${mode} show rank`}
          />
        </label>
        <label className="text-muted flex items-center gap-2 text-xs">
          Stats {!isOverridden(override, "show_stats") && `(${def.show_stats ? "on" : "off"})`}
          <Select
            value={selectValue(override, "show_stats")}
            onValueChange={(v) => setField("show_stats", v)}
            options={OPTIONS}
            aria-label={`${mode} show stats`}
          />
        </label>
      </div>
    </div>
  );
}
