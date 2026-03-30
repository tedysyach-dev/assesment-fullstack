import { cn } from "@/lib/cn";
import * as React from "react";

function Input({ className, type, ...props }: React.ComponentProps<"input">) {
  return (
    <input
      type={type}
      data-slot="input"
      className={cn(
        "px-2 py-3 leading-3.25 tracking-normal border rounded-md text-[14px] border-neutral-50 placeholder:text-neutral-70",
        className,
      )}
      {...props}
    />
  );
}

export { Input };
