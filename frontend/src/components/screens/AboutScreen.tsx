import { useEffect, useState } from "react";
import DOMPurify from "dompurify";
import {
  GetChangelog,
  GetVersion,
} from "../../../bindings/github.com/its-haze/league-rpc/cmd/league-rpc-gui/guiservice";
import { useCheckForUpdates } from "../../hooks/useCheckForUpdates";
import { renderMarkdown } from "../../lib/markdown";
import { Button } from "../ui";
import UpdateBanner from "../UpdateBanner";

// The About section: version, a manual update check independent of whether
// one is already known to be available, and the latest release's changelog.
export function AboutScreen() {
  const [version, setVersion] = useState("");
  const [changelog, setChangelog] = useState<string | null>(null);
  const { checking, result: checkResult, check: handleCheck } = useCheckForUpdates();

  useEffect(() => {
    GetVersion().then(setVersion).catch(() => {});
    GetChangelog().then(setChangelog).catch(() => setChangelog("changelog unavailable"));
  }, []);

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-xl font-semibold">About</h1>

      <section className="border-border bg-surface flex items-center justify-between rounded-lg border p-6">
        <span className="font-mono text-sm">{version ? `v${version}` : "Loading version…"}</span>
        <div className="flex items-center gap-3">
          <Button variant="secondary" onClick={handleCheck} disabled={checking}>
            {checking ? "Checking…" : "Check for updates"}
          </Button>
          {checkResult && <span className="text-muted text-sm">{checkResult}</span>}
        </div>
      </section>

      <UpdateBanner />

      <section className="border-border bg-surface rounded-lg border p-6">
        <h2 className="mb-2 text-sm font-semibold">Changelog</h2>
        <div
          className="prose prose-sm max-w-none text-sm"
          dangerouslySetInnerHTML={{
            __html: DOMPurify.sanitize(renderMarkdown(changelog ?? "Loading changelog…")),
          }}
        />
      </section>
    </div>
  );
}
