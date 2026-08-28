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
// accent rule, plain sentence case.
export function Sidebar({ active, onNavigate }: SidebarProps) {
  return (
    <nav className="border-border bg-surface flex w-44 shrink-0 flex-col gap-1 border-r p-3">
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
    </nav>
  );
}
