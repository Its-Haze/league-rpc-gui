import type { ReactNode } from "react";
import { DiscordIcon, GitHubIcon } from "../icons";
import { DISCORD_COMMUNITY_URL, GITHUB_REPO_URL, openExternal } from "../../lib/links";
import { SECTIONS, type Section } from "../../lib/route";
import { THEME_OPTIONS, type ThemeSetting } from "../../lib/theme";
import { Select } from "../ui";

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
  theme: string;
  onThemeChange: (theme: ThemeSetting) => void;
  /** True until the initial settings load resolves, so the picker can't fire on a config that isn't there yet. */
  themeDisabled?: boolean;
}

// Left nav across the six sections, plus the theme picker: the one setting
// worth reaching from anywhere, so it lives outside any single screen.
export function Sidebar({ active, onNavigate, theme, onThemeChange, themeDisabled }: SidebarProps) {
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

      <div className="border-border flex flex-col gap-2 border-t pt-2">
        <SidebarLink href={DISCORD_COMMUNITY_URL} icon={<DiscordIcon className="size-4" />}>
          Discord
        </SidebarLink>
        <SidebarLink href={GITHUB_REPO_URL} icon={<GitHubIcon className="size-4" />}>
          GitHub
        </SidebarLink>
        <Select
          value={theme}
          onValueChange={(v) => onThemeChange(v as ThemeSetting)}
          options={THEME_OPTIONS}
          disabled={themeDisabled}
          aria-label="Theme"
        />
      </div>
    </nav>
  );
}

function SidebarLink({ href, icon, children }: { href: string; icon: ReactNode; children: ReactNode }) {
  return (
    <a
      href={href}
      onClick={(e) => {
        e.preventDefault();
        openExternal(href);
      }}
      className="text-muted hover:bg-surface-raised hover:text-text flex items-center gap-2 rounded-sm px-3 py-2 text-sm font-medium transition-colors"
    >
      {icon}
      {children}
    </a>
  );
}
