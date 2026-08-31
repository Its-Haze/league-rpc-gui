import { CheckCircle2 } from "lucide-react";

// Onboarding's last screen: a clean sign-off, not another ask. GitHub/Discord already had their moment on the welcome screen.
export function ClosingStep() {
  return (
    <div className="flex flex-col items-center gap-5 text-center">
      <CheckCircle2 className="text-accent size-14" />
      <h1 className="text-3xl font-semibold">You're all set</h1>
      <section className="border-border bg-surface w-full rounded-lg border p-6">
        <p className="text-base">
          League RPC is already running in the background. Closing this window doesn't stop it. It
          keeps going from the tray.
        </p>
      </section>
      <p className="text-muted text-xs">
        Built and maintained by haze. Thanks for giving it a try, and enjoy your games.
      </p>
    </div>
  );
}
