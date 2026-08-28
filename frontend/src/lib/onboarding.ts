// The first-run walkthrough's steps, in display order. Rendered inline, not
// as a modal; the user can revisit an earlier step at any time.
export const ONBOARDING_STEPS = ["welcome", "connect", "display", "startup"] as const;

export type OnboardingStep = (typeof ONBOARDING_STEPS)[number];

export interface OnboardingState {
  step: OnboardingStep;
}

export type OnboardingAction = { type: "next" } | { type: "back" } | { type: "goto"; step: OnboardingStep };

export const initialOnboarding: OnboardingState = { step: ONBOARDING_STEPS[0] };

function indexOf(step: OnboardingStep): number {
  return ONBOARDING_STEPS.indexOf(step);
}

// Pure reducer so step progression is unit-testable without mounting React.
export function onboardingReducer(state: OnboardingState, action: OnboardingAction): OnboardingState {
  switch (action.type) {
    case "next": {
      const i = indexOf(state.step);
      if (i >= ONBOARDING_STEPS.length - 1) return state;
      return { step: ONBOARDING_STEPS[i + 1] };
    }
    case "back": {
      const i = indexOf(state.step);
      if (i <= 0) return state;
      return { step: ONBOARDING_STEPS[i - 1] };
    }
    case "goto":
      return state.step === action.step ? state : { step: action.step };
    default:
      return state;
  }
}

export function isLastStep(step: OnboardingStep): boolean {
  return indexOf(step) === ONBOARDING_STEPS.length - 1;
}

export function isFirstStep(step: OnboardingStep): boolean {
  return indexOf(step) === 0;
}

// The "connect" step's two live checkmarks: both need to be true before it
// makes sense to nudge the user toward "next".
export function connectStepReady(discordConnected: boolean, leagueDetected: boolean): boolean {
  return discordConnected && leagueDetected;
}
