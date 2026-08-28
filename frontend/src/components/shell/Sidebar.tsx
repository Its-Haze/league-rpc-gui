import type { ReactNode } from "react";
import { DiscordIcon, GitHubIcon } from "../icons";
import { DISCORD_COMMUNITY_URL, GITHUB_REPO_URL } from "../../lib/links";
import { SECTIONS, type Section } from "../../lib/route";

const LABELS: Record<Section, string> = {
  home: "Home",
  display: "Display",
  behavior: "Behavior",
  advanced: "Advanced",
  help: "Help",
  about: "About",
};

export interface SidebarProps {
  active: Section;
  onNavigate: (section: Section) => void;
}

// Left nav across the six sections. The active item is a filled row, no
export function Sidebar({ active, onNavigate }: SidebarProps) {
  return (
    <nav className="border-border bg-surface flex w-44 shrink-0 flex-col border-r p-3">
      <div className="flex flex-1 flex-col gap-1">
        {SECTIONS.map((section) => {
          const isActive = section === active;
          return (
            <button
              key={section}
              onClick={() => onNavigate(section)}
              aria-current={isActive ? "page" : undefined}
              className={
                "rounded-sm px-3 py-2 text-left text-sm font-medium transition-colors " +
                (isActive
                  ? "bg-surface-raised text-accent"
                  : "text-muted hover:bg-surface-raised")
              }
            >
              {LABELS[section]}
            </button>
          );
        })}
      </div>

      <div className="border-border flex flex-col gap-1 border-t pt-2">
        <SidebarLink href={DISCORD_COMMUNITY_URL} icon={<DiscordIcon className="size-4" />}>
          Discord
        </SidebarLink>
        <SidebarLink href={GITHUB_REPO_URL} icon={<GitHubIcon className="size-4" />}>
          GitHub
        </SidebarLink>
      </div>
    </nav>
  );
}

function SidebarLink({ href, icon, children }: { href: string; icon: ReactNode; children: ReactNode }) {
  return (
    <a
      href={href}
      target="_blank"
      rel="noreferrer"
      className="text-muted hover:bg-surface-raised hover:text-text flex items-center gap-2 rounded-sm px-3 py-2 text-sm font-medium transition-colors"
    >
      {icon}
      {children}
    </a>
  );
}
