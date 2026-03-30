// @/features/dashboard/components/DashboardLayout.tsx
import { Outlet, NavLink } from "react-router-dom";
import { BellIcon, LogOutIcon, MenuIcon, XIcon } from "lucide-react";
import { AppsIcon } from "@/assets/icon";
import { useAuthStore } from "@/features/auth/store/authStore";
import { useState } from "react";

const NAV_ITEMS = [
  { label: "Inbound", to: "/dashboard/inbound" },
  { label: "Outbound", to: "/dashboard/outbound" },
  { label: "Inventory", to: "/dashboard/inventory" },
  { label: "Settings", to: "/dashboard/settings" },
];

const DashboardLayout = () => {
  const logout = useAuthStore((state) => state.logout);
  const [open, setOpen] = useState(false);

  return (
    <div className="flex h-screen">
      <div className="flex flex-col flex-1 overflow-hidden">
        <nav className="w-full py-4 px-4 md:px-10 bg-primary-main shrink-0 flex items-center justify-between relative">
          <div className="flex items-center gap-3">
            <div className="bg-white rounded-md p-2 w-9 h-9 flex items-center justify-center shadow-sm">
              <AppsIcon />
            </div>

            <span className="text-white text-sm md:text-base tracking-widest">
              WMSpaceIO
            </span>
          </div>

          <div className="hidden md:flex items-center gap-1">
            {NAV_ITEMS.map(({ label, to }) => (
              <NavLink
                key={to}
                to={to}
                className={({ isActive }) =>
                  `flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-colors
                  ${
                    isActive
                      ? "bg-white text-primary-main"
                      : "text-white/70 hover:bg-white/10 hover:text-white"
                  }`
                }
              >
                {label}
              </NavLink>
            ))}
          </div>

          <div className="flex items-center gap-2 md:gap-3">
            <button className="relative bg-white p-2 rounded-lg">
              <BellIcon className="w-5 h-5 transform rotate-45 hover:rotate-0 transition-transform" />

              <span className="absolute -top-1 -right-1 min-w-4 h-4 px-1 bg-red-500 text-white text-[10px] font-semibold rounded-full flex items-center justify-center">
                2
              </span>
            </button>

            <div className="w-8 h-8 rounded-full bg-white/20 border-2 border-white/30 flex items-center justify-center">
              <span className="text-white text-xs font-semibold">A</span>
            </div>

            <button
              onClick={logout}
              className="hidden md:block text-white/70 hover:text-white p-2 rounded-lg hover:bg-white/10"
            >
              <LogOutIcon className="w-5 h-5" />
            </button>

            <button
              className="md:hidden text-white p-2"
              onClick={() => setOpen(!open)}
            >
              {open ? <XIcon /> : <MenuIcon />}
            </button>
          </div>

          {open && (
            <div className="absolute top-full left-0 w-full bg-primary-main border-t border-white/10 flex flex-col md:hidden p-2">
              {NAV_ITEMS.map(({ label, to }) => (
                <NavLink
                  key={to}
                  to={to}
                  onClick={() => setOpen(false)}
                  className={({ isActive }) =>
                    `flex items-center gap-2 px-4 py-3 rounded-lg text-sm font-medium transition-colors
                    ${
                      isActive
                        ? "bg-white text-primary-main"
                        : "text-white/70 hover:bg-white/10 hover:text-white"
                    }`
                  }
                >
                  {label}
                </NavLink>
              ))}

              <button
                onClick={logout}
                className="flex items-center gap-2 px-4 py-3 text-white/70 hover:text-white"
              >
                <LogOutIcon className="w-4 h-4" />
                Logout
              </button>
            </div>
          )}
        </nav>

        <main className="flex-1 overflow-y-auto p-4 md:p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default DashboardLayout;
