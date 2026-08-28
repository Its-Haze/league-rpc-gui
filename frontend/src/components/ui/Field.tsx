import type { ReactNode } from "react";
import * as LabelPrimitive from "@radix-ui/react-label";

export interface FieldProps {
  id: string;
  label: string;
  hint?: string;
  children: ReactNode;
}

// Label + control + optional hint, the row shape every settings screen
export function Field({ id, label, hint, children }: FieldProps) {
  return (
    <div className="flex items-center justify-between gap-4 py-2">
      <div className="flex flex-col gap-0.5">
        <LabelPrimitive.Root htmlFor={id} className="text-sm">
          {label}
        </LabelPrimitive.Root>
        {hint && <span className="text-muted text-xs">{hint}</span>}
      </div>
      {children}
    </div>
  );
}
