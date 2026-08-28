import { useEffect, useRef, useState } from "react";
import { Events } from "@wailsio/runtime";
import {
  GetSettings,
  ApplySettings,
  GetPresets,
  GetVersion,
} from "../bindings/github.com/its-haze/league-rpc/cmd/league-rpc-gui/guiservice";
import type { Config } from "../bindings/github.com/its-haze/league-rpc/internal/config/models";
import { AppShell } from "./components/shell/AppShell";
import UpdateBanner from "./components/UpdateBanner";
import { useAppliedTheme } from "./hooks/useAppliedTheme";
import type { ThemeSetting } from "./lib/theme";

const CONFIG_CHANGED_EVENT = "settings:changed";

export default function App() {
  const [cfg, setCfg] = useState<Config | null>(null);
  const [presets, setPresets] = useState<Record<string, string | undefined>>({});
  const [version, setVersion] = useState<string>("");
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  useAppliedTheme(cfg?.theme ?? "system");

  // Mirrors cfg synchronously so back-to-back applyPatch calls build their
  // patch on top of each other's optimistic write, not a stale render closure.
  const cfgRef = useRef<Config | null>(null);
  cfgRef.current = cfg;

  useEffect(() => {
    GetSettings().then(setCfg).catch((e: unknown) => setError(String(e)));
    GetPresets()
      .then((p) => setPresets(p ?? {}))
      .catch(() => {});
    GetVersion().then(setVersion).catch(() => {});

    const off = Events.On(CONFIG_CHANGED_EVENT, (ev: { data: Config }) => {
      setCfg(ev.data);
      cfgRef.current = ev.data;
    });
    return () => off();
  }, []);

  async function applyPatch(patch: Partial<Config>) {
    const current = cfgRef.current;
    if (!current) return;
    const next: Config = { ...current, ...patch };
    cfgRef.current = next;
    setCfg(next);
    setSaving(true);
    setError(null);
    try {
      await ApplySettings(next);
    } catch (e) {
      setError(String(e));
      cfgRef.current = current;
      setCfg(current);
    } finally {
      setSaving(false);
    }
  }

  async function toggleEmojis() {
    if (!cfg) return;
    await applyPatch({ presence: { ...cfg.presence, show_emojis: !cfg.presence.show_emojis } });
  }

  function handleThemeChange(theme: ThemeSetting) {
    void applyPatch({ theme });
  }

  return (
    <AppShell
      theme={cfg?.theme ?? "system"}
      onThemeChange={handleThemeChange}
      themeDisabled={!cfg}
      error={error}
      homeContent={
        <>
          {!cfg ? (
            <p className="text-muted text-sm">Loading settings…</p>
          ) : (
            <section className="border-border bg-surface flex flex-col gap-4 rounded-lg border p-6">
              <Row label="Schema version" value={String(cfg.schema_version)} />
              <Row label="Discord App ID" value={cfg.discord_app_id} />
              <Row label="Presets" value={Object.keys(presets).join(", ") || "none"} />

              <label className="flex items-center justify-between gap-4">
                <span className="text-sm">Show status emojis in presence</span>
                <button
                  onClick={toggleEmojis}
                  disabled={saving}
                  className="border-border bg-surface-raised rounded-sm border px-3 py-1 text-sm"
                >
                  {cfg.presence.show_emojis ? "On" : "Off"}
                </button>
              </label>
            </section>
          )}
        </>
      }
      aboutContent={
        <>
          <p className="text-muted text-sm">
            {version ? <span className="font-mono">v{version}</span> : "Loading version…"}
          </p>
          <UpdateBanner />
        </>
      }
    />
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between gap-4 text-sm">
      <span className="text-muted">{label}</span>
      <span className="font-mono">{value}</span>
    </div>
  );
}
