import { useEffect, useState } from "react";
import { Events } from "@wailsio/runtime";
import DOMPurify from "dompurify";
import {
  CheckForUpdates,
  GetChangelog,
  GetUpdateStatus,
  RestartForUpdate,
  StartUpdate,
} from "../../bindings/github.com/its-haze/league-rpc/cmd/league-rpc-gui/guiservice";
import type { UpdateStatus } from "../../bindings/github.com/its-haze/league-rpc/internal/app/models";
import { renderMarkdown } from "../lib/markdown";

// Event names the Wails updater emits directly on the same event bus as our
// own "update:changed" — see internal/updates and pkg/updater/events.go.
const EVENT_CHANGED = "update:changed";
const EVENT_DOWNLOAD_PROGRESS = "wails:updater:download-progress";
const EVENT_UPDATE_READY = "wails:updater:update-ready";
const EVENT_ERROR = "wails:updater:error";

type Phase = "idle" | "downloading" | "ready" | "error";

interface Progress {
  written: number;
  total: number;
}

export default function UpdateBanner() {
  const [status, setStatus] = useState<UpdateStatus | null>(null);
  const [dismissedVersion, setDismissedVersion] = useState<string | null>(null);
  const [phase, setPhase] = useState<Phase>("idle");
  const [progress, setProgress] = useState<Progress | null>(null);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [changelog, setChangelog] = useState<string | null>(null);
  const [showChangelog, setShowChangelog] = useState(false);
  const [checking, setChecking] = useState(false);

  useEffect(() => {
    GetUpdateStatus().then(setStatus).catch(() => {});

    const offs = [
      Events.On(EVENT_CHANGED, (ev: { data: UpdateStatus }) => setStatus(ev.data)),
      Events.On(EVENT_DOWNLOAD_PROGRESS, (ev: { data: Progress }) => {
        setPhase("downloading");
        setProgress(ev.data);
      }),
      Events.On(EVENT_UPDATE_READY, () => {
        setPhase("ready");
        setProgress(null);
      }),
      Events.On(EVENT_ERROR, (ev: { data: { message?: string } }) => {
        setPhase("error");
        setErrorMsg(ev.data?.message ?? "the update failed");
      }),
    ];
    return () => offs.forEach((off) => off());
  }, []);

  if (!status?.available || status.version === dismissedVersion) {
    return null;
  }

  async function handleCheck() {
    setChecking(true);
    try {
      await CheckForUpdates();
    } finally {
      setChecking(false);
    }
  }

  async function handleUpdate() {
    setPhase("downloading");
    setErrorMsg(null);
    try {
      await StartUpdate();
      setPhase("ready");
    } catch (e) {
      setPhase("error");
      setErrorMsg(String(e));
    }
  }

  async function handleRestart() {
    try {
      await RestartForUpdate();
    } catch (e) {
      setPhase("error");
      setErrorMsg(String(e));
    }
  }

  async function handleToggleChangelog() {
    const next = !showChangelog;
    setShowChangelog(next);
    if (next && changelog === null) {
      setChangelog(await GetChangelog());
    }
  }

  return (
    <div className="border-accent bg-surface-raised flex flex-col gap-3 rounded-lg border p-4 text-sm">
      <div className="flex items-center justify-between gap-4">
        <div>
          <strong>League RPC {status.version}</strong> is available.
          {phase === "downloading" && progress && (
            <span className="text-muted ml-2">
              downloading{progressPercent(progress)}
            </span>
          )}
          {phase === "ready" && <span className="text-ok ml-2">ready to install</span>}
          {phase === "error" && (
            <span className="text-danger ml-2">{errorMsg ?? "update failed"}</span>
          )}
        </div>
        <div className="flex items-center gap-2">
          <button onClick={handleToggleChangelog} className="text-muted underline">
            {showChangelog ? "Hide changelog" : "Changelog"}
          </button>
          {phase === "ready" ? (
            <button
              onClick={handleRestart}
              className="bg-accent text-accent-text rounded-sm px-3 py-1"
            >
              Restart now
            </button>
          ) : (
            <button
              onClick={handleUpdate}
              disabled={phase === "downloading"}
              className="bg-accent text-accent-text rounded-sm px-3 py-1 disabled:opacity-60"
            >
              {phase === "downloading" ? "Updating…" : "Update"}
            </button>
          )}
          <button
            onClick={() => setDismissedVersion(status.version)}
            aria-label="Dismiss"
            className="text-muted px-1"
          >
            ×
          </button>
        </div>
      </div>

      {showChangelog && (
        <div
          className="border-border prose prose-sm max-w-none border-t pt-3"
          dangerouslySetInnerHTML={{
            __html: DOMPurify.sanitize(renderMarkdown(changelog ?? "Loading changelog…")),
          }}
        />
      )}

      <div className="text-muted flex items-center gap-3">
        <button onClick={handleCheck} disabled={checking} className="underline disabled:opacity-60">
          {checking ? "Checking…" : "Check for updates"}
        </button>
      </div>
    </div>
  );
}

function progressPercent(p: Progress): string {
  if (p.total <= 0) return "";
  return ` ${Math.round((p.written / p.total) * 100)}%`;
}
