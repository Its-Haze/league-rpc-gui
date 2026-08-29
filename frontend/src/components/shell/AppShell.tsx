import type { ReactNode } from "react";
import { useCloseRequest } from "../../hooks/useCloseRequest";
import { useRoute } from "../../hooks/useRoute";
import type { ThemeSetting } from "../../lib/theme";
import { CloseConfirmDialog } from "./CloseConfirmDialog";
import { Sidebar } from "./Sidebar";
import { TitleBar } from "./TitleBar";
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
}

// The app shell: title bar, sidebar, and screen content. Sidebar stays
// mounted across every section, so the theme picker never remounts on navigation.
export function AppShell({ theme, onThemeChange, themeDisabled, error }: AppShellProps) {
  const { section, navigate } = useRoute();
  const closeRequest = useCloseRequest();

  return (
    <div className="flex h-full flex-col">
      <CloseConfirmDialog open={closeRequest.pending} onDismiss={closeRequest.dismiss} />
      <TitleBar />
      <div className="flex min-h-0 flex-1">
        <Sidebar
          active={section}
          onNavigate={navigate}
          theme={theme}
          onThemeChange={onThemeChange}
          themeDisabled={themeDisabled}
        />
        <div className="flex min-w-0 flex-1 flex-col">
          {error && (
            <div className="border-danger text-danger border-b px-6 py-2 text-sm">{error}</div>
          )}
          <main className="min-w-0 flex-1 overflow-y-auto p-6">
            {section === "home" && <HomeScreen />}
            {section === "display" && <DisplayScreen />}
            {section === "behavior" && <BehaviorScreen />}
            {section === "advanced" && <AdvancedScreen />}
            {section === "help" && <HelpScreen />}
            {section === "about" && <AboutScreen />}
          </main>
        </div>
      </div>
    </div>
  );
}
