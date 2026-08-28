import type { ReactNode } from "react";

// Placeholder host for the Home dashboard. The onboarding/status content
// itself lands in ticket 15; for now it hosts what App.tsx already had.
export function HomeScreen({ children }: { children?: ReactNode }) {
  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-xl font-semibold">Home</h1>
      {children}
    </div>
  );
}
