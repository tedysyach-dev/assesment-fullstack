import api from "@/lib/axios";
import { useAuthStore } from "../store/authStore";
import type { LoginRequest } from "../types";

export const useAuth = () => {
  const { isAuthenticated, setAuth, logout } = useAuthStore();

  const login = async (payload: LoginRequest) => {
    try {
      const { data } = await api.post("/auth/login", payload);

      if (!data.status || !data.resource?.accessToken) {
        throw new Error(data.message ?? "Login failed");
      }

      setAuth(data.resource.accessToken);
    } catch (error: any) {
      const message =
        error?.response?.data?.message ?? error?.message ?? "Login failed";

      throw new Error(message);
    }
  };

  return { isAuthenticated, login, logout };
};
