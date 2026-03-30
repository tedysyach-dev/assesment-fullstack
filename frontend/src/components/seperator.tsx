import * as React from "react";
import { cn } from "@/lib/cn";

function Separator({
  className,
  orientation = "horizontal",
  decorative = true,
  ...props
}: React.ComponentProps<"div"> & {
  orientation?: "horizontal" | "vertical";
  decorative?: boolean;
}) {
  return (
    <div
      data-slot="separator"
      role={decorative ? "none" : "separator"}
      aria-orientation={!decorative ? orientation : undefined}
      data-orientation={orientation}
      className={cn(
        "bg-neutral-50 shrink-0",
        orientation === "horizontal" ? "h-px w-full" : "w-px self-stretch",
        className,
      )}
      {...props}
    />
  );
}

function FieldSeparator({
  children,
  className,
  ...props
}: React.ComponentProps<"div"> & {
  children?: React.ReactNode;
}) {
  return (
    <div
      data-slot="field-separator"
      className={cn("flex items-center gap-3 w-full", className)}
      {...props}
    >
      <Separator className="flex-1" />
      {children && (
        <span className="text-neutral-80 text-xs whitespace-nowrap">
          {children}
        </span>
      )}
      <Separator className="flex-1" />
    </div>
  );
}

export { Separator, FieldSeparator };
