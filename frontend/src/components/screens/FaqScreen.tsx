import { FAQ_GROUPS, type FaqEntry } from "../../lib/faq";
import { openExternal } from "../../lib/links";
import { SettingsCard } from "../ui";

// The FAQ section: one card per topic in a two-column grid, every answer
// visible at once. No accordion, so nothing is a click away and Ctrl+F works.
export function FaqScreen() {
  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-xl font-semibold">FAQ</h1>

      <div className="grid gap-4 md:grid-cols-2">
        {FAQ_GROUPS.map((group) => (
          <SettingsCard key={group.id} icon={group.icon} title={group.title}>
            <div className="flex flex-col gap-4 pt-1">
              {group.entries.map((entry) => (
                <Entry key={entry.question} entry={entry} />
              ))}
            </div>
          </SettingsCard>
        ))}
      </div>
    </div>
  );
}

function Entry({ entry }: { entry: FaqEntry }) {
  return (
    <div className="flex flex-col gap-1">
      <h3 className="text-sm font-semibold">{entry.question}</h3>
      <p className="text-muted text-sm leading-relaxed">{entry.answer}</p>
      {entry.links && (
        <div className="flex flex-wrap gap-4 pt-0.5">
          {entry.links.map((link) => (
            <a
              key={link.href}
              href={link.href}
              onClick={(e) => {
                e.preventDefault();
                openExternal(link.href);
              }}
              className="text-accent text-sm font-medium hover:underline"
            >
              {link.label}
            </a>
          ))}
        </div>
      )}
    </div>
  );
}
