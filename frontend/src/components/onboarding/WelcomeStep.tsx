import { DiscordIcon, GitHubIcon } from "../icons";
import { DISCORD_COMMUNITY_URL, GITHUB_REPO_URL, openExternal } from "../../lib/links";
import { Button } from "../ui";

// Onboarding's first screen: what the app does, and the same GitHub/Discord
// CTAs repeated on the closing screen.
export function WelcomeStep() {
  return (
    <div className="flex flex-col gap-4 text-center">
      <h1 className="text-2xl font-semibold">Welcome to League RPC</h1>
      <p className="text-muted text-sm">
        League RPC turns what you're doing in League of Legends into your Discord status, updated
        automatically while you play.
      </p>
      <ul className="border-border bg-surface flex list-disc flex-col gap-2 rounded-lg border p-4 pl-8 text-left text-sm">
        <li>Shows your queue, champion, and rank the moment you load into a game.</li>
        <li>Tracks live KDA and creep score without you touching a thing.</li>
        <li>Runs from the tray: start it once and forget about it.</li>
        <li>Everything here is yours to customize, starting on the next screen.</li>
      </ul>
      <div className="flex items-center justify-center gap-3">
        <Button variant="secondary" asChild>
          <a
            href={GITHUB_REPO_URL}
            onClick={(e) => {
              e.preventDefault();
              openExternal(GITHUB_REPO_URL);
            }}
          >
            <GitHubIcon className="size-4" />
            Star on GitHub
          </a>
        </Button>
        <Button variant="secondary" asChild>
          <a
            href={DISCORD_COMMUNITY_URL}
            onClick={(e) => {
              e.preventDefault();
              openExternal(DISCORD_COMMUNITY_URL);
            }}
          >
            <DiscordIcon className="size-4" />
            Join the Discord
          </a>
        </Button>
      </div>
    </div>
  );
}
