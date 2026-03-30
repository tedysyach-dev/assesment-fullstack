import { createBrowserRouter, RouterProvider } from "react-router-dom";
import ProtectedRoute from "@/features/auth/components/ProtectedRoute";
import LoginPage from "@/features/auth/components/LoginPage";
import OrderPage from "@/features/dashboard/outbound/components/OrderPage";
import DashboardLayout from "@/features/dashboard/components/DashboardLayout";

const router = createBrowserRouter([
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    // ini adalah route dashboard
    element: <ProtectedRoute />,
    path: "/dashboard",
    children: [
      {
        element: <DashboardLayout />,
        children: [
          {
            path: "outbound",
            element: <OrderPage />,
          },
        ],
      },
    ],
  },
]);

const AppRouter = () => <RouterProvider router={router} />;

export default AppRouter;
