import type { ReactNode } from "react";

// Placeholder host for the About screen. Links, diagnostics, and the full
// changelog view land in ticket 16; for now it hosts the version/update banner.
export function AboutScreen({ children }: { children?: ReactNode }) {
  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-xl font-semibold">About</h1>
      {children}
    </div>
  );
}
