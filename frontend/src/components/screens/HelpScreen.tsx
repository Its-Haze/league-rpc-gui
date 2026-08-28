import { useState } from "react";
import {
  GetDiagnostics,
  OpenLogsFolder,
} from "../../../bindings/github.com/its-haze/league-rpc/cmd/league-rpc-gui/guiservice";
import { useSettings } from "../../hooks/useSettings";
import { formatDiagnostics } from "../../lib/diagnostics";
import { BUG_REPORT_URL, DISCORD_COMMUNITY_URL, FEATURE_REQUEST_URL } from "../../lib/links";
import { Button } from "../ui";
import { LogViewer } from "./help/LogViewer";

// The Help section: log viewer, logs folder, diagnostics, and where to go
// for support or to file an issue.
export function HelpScreen() {
  const { cfg, applyPatch } = useSettings();
  const [copyResult, setCopyResult] = useState<string | null>(null);
  const [folderError, setFolderError] = useState<string | null>(null);

  async function handleOpenFolder() {
    setFolderError(null);
    try {
      await OpenLogsFolder();
    } catch (e) {
      setFolderError(String(e));
    }
  }

  async function handleCopyDiagnostics() {
    try {
      const diag = await GetDiagnostics();
      await navigator.clipboard.writeText(formatDiagnostics(diag));
      setCopyResult("Copied to clipboard.");
    } catch (e) {
      setCopyResult(String(e));
    }
  }

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-xl font-semibold">Help</h1>

      <section className="border-border bg-surface flex flex-col gap-3 rounded-lg border p-6">
        <div className="flex items-center justify-between">
          <h2 className="text-sm font-semibold">Logs</h2>
          <div className="flex items-center gap-2">
            <Button variant="secondary" onClick={handleOpenFolder}>
              Open logs folder
            </Button>
            <Button variant="secondary" onClick={handleCopyDiagnostics}>
              Copy diagnostics
            </Button>
          </div>
        </div>
        <LogViewer />
        {folderError && <p className="text-danger text-xs">{folderError}</p>}
        {copyResult && <p className="text-muted text-xs">{copyResult}</p>}
      </section>

      <section className="border-border bg-surface flex flex-col gap-2 rounded-lg border p-6">
        <h2 className="text-sm font-semibold">Get help</h2>
        <a href={DISCORD_COMMUNITY_URL} target="_blank" rel="noreferrer" className="text-accent text-sm underline">
          Join the Discord community
        </a>
        <a href={BUG_REPORT_URL} target="_blank" rel="noreferrer" className="text-accent text-sm underline">
          Report a bug
        </a>
        <a href={FEATURE_REQUEST_URL} target="_blank" rel="noreferrer" className="text-accent text-sm underline">
          Request a feature
        </a>
      </section>

      {cfg && (
        <section className="border-border bg-surface flex items-center justify-between rounded-lg border p-6">
          <div>
            <h2 className="text-sm font-semibold">First-run walkthrough</h2>
            <p className="text-muted text-xs">Replay the setup steps shown on first launch.</p>
          </div>
          <Button variant="secondary" onClick={() => void applyPatch({ onboarding_complete: false })}>
            Replay walkthrough
          </Button>
        </section>
      )}
    </div>
  );
}
