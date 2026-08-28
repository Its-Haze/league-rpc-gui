import type { ReactNode } from "react";
import { useRoute } from "../../hooks/useRoute";
import type { ThemeSetting } from "../../lib/theme";
import { Sidebar } from "./Sidebar";
import { TopStrip } from "./TopStrip";
import { HomeScreen } from "../screens/HomeScreen";
import { DisplayScreen } from "../screens/DisplayScreen";
import { BehaviorScreen } from "../screens/BehaviorScreen";
import { AdvancedScreen } from "../screens/AdvancedScreen";
import { HelpScreen } from "../screens/HelpScreen";
import { AboutScreen } from "../screens/AboutScreen";

export interface AppShellProps {
  theme: string;
  onThemeChange: (theme: ThemeSetting) => void;
  themeDisabled?: boolean;
  error?: ReactNode;
  homeContent: ReactNode;
  aboutContent: ReactNode;
}

// The app shell: sidebar + top strip stay mounted across every section,
// so the connection lights and Pause toggle never remount on navigation.
export function AppShell({
  theme,
  onThemeChange,
  themeDisabled,
  error,
  homeContent,
  aboutContent,
}: AppShellProps) {
  const { section, navigate } = useRoute();

  return (
    <div className="flex h-full">
      <Sidebar active={section} onNavigate={navigate} />
      <div className="flex min-w-0 flex-1 flex-col">
        <TopStrip theme={theme} onThemeChange={onThemeChange} themeDisabled={themeDisabled} />
        {error && (
          <div className="border-danger text-danger border-b px-6 py-2 text-sm">{error}</div>
        )}
        <main className="min-w-0 flex-1 overflow-y-auto p-6">
          {section === "home" && <HomeScreen>{homeContent}</HomeScreen>}
          {section === "display" && <DisplayScreen />}
          {section === "behavior" && <BehaviorScreen />}
          {section === "advanced" && <AdvancedScreen />}
          {section === "help" && <HelpScreen />}
          {section === "about" && <AboutScreen>{aboutContent}</AboutScreen>}
        </main>
      </div>
    </div>
  );
}
