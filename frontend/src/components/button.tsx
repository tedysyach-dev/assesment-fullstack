import { cn } from "@/lib/cn";
import { cva, type VariantProps } from "class-variance-authority";
import * as React from "react";

const buttonVariants = cva(
  [
    "inline-flex items-center justify-center gap-1.5",
    "rounded-lg border border-transparent",
    "text-sm font-medium whitespace-nowrap select-none",
    "cursor-pointer transition-all outline-none shrink-0",
    "disabled:pointer-events-none disabled:opacity-50",
    "focus-visible:ring-3 focus-visible:ring-primary-main/50",
    "[&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
  ],
  {
    variants: {
      variant: {
        default: "bg-primary-main !text-white hover:bg-primary-pressed",
        outline:
          "border-primary-border text-neutral-100 bg-transparent hover:bg-primary-surface",
        ghost: "text-primary-main bg-transparent hover:bg-primary-surface",
        destructive: "bg-danger-main !text-white hover:bg-danger-main/80",
      },
      size: {
        sm: "h-7 px-2 text-xs",
        default: " w-full py-2.5 text-xs",
        lg: "h-[42px] px-4 text-base",
        icon: "size-8 p-0",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  },
);

function Button({
  className,
  variant,
  size,
  type = "button", // tambah ini
  ...props
}: React.ComponentProps<"button"> & VariantProps<typeof buttonVariants>) {
  return (
    <button
      data-slot="button"
      type={type} // tambah ini
      className={cn(buttonVariants({ variant, size, className }))}
      {...props}
    />
  );
}

export { Button, buttonVariants };
