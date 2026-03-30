import { Check, Loader2Icon, X } from "lucide-react";
import { Toaster as Sonner, type ToasterProps } from "sonner";

const Toaster = ({ style, ...props }: ToasterProps) => {
  return (
    <Sonner
      theme="system"
      className="toaster group"
      closeButton
      icons={{
        success: (
          <span className="flex items-center rounded-full bg-semantic-green/30 justify-center size-6">
            <Check className="size-3 text-semantic-green" />
          </span>
        ),
        error: (
          <span className="flex items-center rounded-full bg-danger-main/30 justify-center size-6">
            <X className="size-3 text-danger-main" />
          </span>
        ),
        loading: <Loader2Icon className="size-4 animate-spin" />,
      }}
      toastOptions={{
        classNames: {
          toast:
            "!bg-[#1e2a2a] !border-none !rounded-[14px] px-[12px] font-poppins",
          title: "!text-white !text-sm !font-semibold",
          description: "!text-neutral-30 !text-xs !font-light w-[244px]",
          closeButton:
            "!bg-transparent !border-none !text-white !right-2 !left-auto !top-4 !size-4 [&>svg]:!size-[30px]",
          icon: "!rounded-full",
          error: "!bg-[linear-gradient(35deg,#E52A3450_0%,#353F3D_35%)]",
          success: "!bg-[linear-gradient(35deg,#4CAF5050_0%,#353F3D_35%)]",
        },
      }}
      style={
        {
          "--normal-bg": "var(--popover)",
          "--normal-text": "var(--popover-foreground)",
          "--normal-border": "var(--border)",
          "--border-radius": "var(--radius)",
          "--toast-icon-margin-end": "12px",
          ...style, // merge disini
        } as React.CSSProperties
      }
      {...props}
    />
  );
};

export { Toaster };
