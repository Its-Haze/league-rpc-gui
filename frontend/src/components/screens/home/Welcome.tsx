// Static, non-technical introduction: what the app does and a few things
// worth knowing before using it. Short, read-once copy, not a manual.
export function Welcome() {
  return (
    <section className="border-border bg-surface flex flex-col gap-3 rounded-lg border p-6">
      <div className="flex flex-col gap-1">
        <h2 className="text-sm font-semibold">Welcome to League RPC</h2>
        <p className="text-sm">
          League RPC turns what you're doing in League of Legends into your Discord status,
          updated automatically while you play.
        </p>
      </div>

      <ul className="flex list-disc flex-col gap-1 pl-5 text-sm">
        <li>Start it whenever: before League, after League, or mid-game. It catches up on its own.</li>
        <li>
          Turn on "Start with Windows" in Behavior so it's always running and you don't have to
          remember to launch it.
        </li>
        <li>
          Closing this window doesn't quit League RPC. It keeps running in the tray; right-click
          the icon and choose Quit to stop it fully.
        </li>
        <li>Every setting here is saved automatically and stays that way between launches.</li>
      </ul>
    </section>
  );
}
