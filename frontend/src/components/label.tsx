import { cn } from "@/lib/cn";
import * as React from "react";

function Label({ className, ...props }: React.ComponentProps<"label">) {
  return (
    <label
      data-slot="label"
      className={cn(
        "gap-2 text-sm font-jakarta text-neutral-100 leading-none font-semibold flex items-center select-none",
        className,
      )}
      {...props}
    />
  );
}

export { Label };
