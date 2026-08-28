import { RotateCcw } from "lucide-react";
import type { ReactNode } from "react";
import * as LabelPrimitive from "@radix-ui/react-label";

export interface FieldProps {
  id: string;
  label: string;
  hint?: string;
  children: ReactNode;
  /** Reverts this setting to its built-in default. Omit when the setting has
   * no meaningful default (e.g. a free-form path) or is reset some other way. */
  onReset?: () => void;
  /** True once the value already matches the default, so the reset control
   * can hide instead of doing nothing when clicked. */
  isDefault?: boolean;
}

// Label + control + optional hint, the row shape every settings screen
export function Field({ id, label, hint, children, onReset, isDefault }: FieldProps) {
  return (
    <div className="flex items-center justify-between gap-4 py-2">
      <div className="flex flex-col gap-0.5">
        <LabelPrimitive.Root htmlFor={id} className="text-sm">
          {label}
        </LabelPrimitive.Root>
        {hint && <span className="text-muted text-xs">{hint}</span>}
      </div>
      <div className="flex items-center gap-2">
        {children}
        {onReset && (
          <button
            type="button"
            onClick={onReset}
            title="Reset to default"
            aria-label={`Reset ${label} to default`}
            disabled={isDefault}
            // Reserved even while at default so the control beside it doesn't
            // shift width when the button appears; just fades out instead.
            className="text-muted hover:text-text shrink-0 rounded-sm p-1 disabled:pointer-events-none disabled:opacity-0"
          >
            <RotateCcw className="size-3.5" />
          </button>
        )}
      </div>
    </div>
  );
}
