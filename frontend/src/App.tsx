import { useEffect, useState } from "react";
import { Events } from "@wailsio/runtime";
import {
  GetSettings,
  ApplySettings,
  GetPresets,
  GetVersion,
} from "../bindings/github.com/its-haze/league-rpc/cmd/league-rpc-gui/guiservice";
import type { Config } from "../bindings/github.com/its-haze/league-rpc/internal/config/models";
import UpdateBanner from "./components/UpdateBanner";

// Placeholder shell. Real screens (Home, Display, Behavior, Advanced) land in
export default function App() {
  const [cfg, setCfg] = useState<Config | null>(null);
  const [presets, setPresets] = useState<Record<string, string | undefined>>({});
  const [version, setVersion] = useState<string>("");
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    GetSettings().then(setCfg).catch((e: unknown) => setError(String(e)));
    GetPresets()
      .then((p) => setPresets(p ?? {}))
      .catch(() => {});
    GetVersion().then(setVersion).catch(() => {});

    const off = Events.On("settings:changed", (ev: { data: Config }) => {
      setCfg(ev.data);
    });
    return () => off();
  }, []);

  async function toggleEmojis() {
    if (!cfg) return;
    const next: Config = {
      ...cfg,
      presence: { ...cfg.presence, show_emojis: !cfg.presence.show_emojis },
    };
    setSaving(true);
    setError(null);
    try {
      await ApplySettings(next);
      setCfg(next);
    } catch (e) {
      setError(String(e));
    } finally {
      setSaving(false);
    }
  }

  return (
    <main className="mx-auto flex min-h-full max-w-2xl flex-col gap-6 p-8">
      <header className="flex flex-col gap-1">
        <h1 className="text-xl font-semibold">League RPC</h1>
        <p className="text-muted text-sm">
          The daemon is running in this process. Settings below are live.
          {version && <span className="font-mono"> · v{version}</span>}
        </p>
      </header>

      <UpdateBanner />

      {error && (
        <div className="border-danger text-danger rounded-md border px-3 py-2 text-sm">
          {error}
        </div>
      )}

      {!cfg ? (
        <p className="text-muted text-sm">Loading settingsâ€¦</p>
      ) : (
        <section className="border-border bg-surface flex flex-col gap-4 rounded-lg border p-6">
          <Row label="Schema version" value={String(cfg.schema_version)} />
          <Row label="Discord App ID" value={cfg.discord_app_id} />
          <Row label="Theme" value={cfg.theme} />
          <Row
            label="Presets"
            value={Object.keys(presets).join(", ") || "none"}
          />

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
    </main>
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
