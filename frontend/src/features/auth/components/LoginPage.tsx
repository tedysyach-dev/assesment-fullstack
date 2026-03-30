import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";
import { Input } from "@/components/input";
import { Label } from "@/components/label";
import { Button } from "@/components/button";
import { FieldSeparator } from "@/components/seperator";
import { GoogleIcon } from "@/assets/google";
import { AppsIcon } from "@/assets/icon";
import { Toaster } from "@/components/toast";
import { toast } from "sonner";

const LoginPage = () => {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [rememberMe, setRememberMe] = useState(false);
  const [emailError, setEmailError] = useState(""); // tambahan

  const validateEmail = (value: string) => {
    if (!value) {
      return "Email tidak boleh kosong";
    }
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(value)) {
      return "Format email tidak valid";
    }
    return "";
  };

  const handleEmailBlur = () => {
    const errorMsg = validateEmail(email);
    setEmailError(errorMsg);
  };

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value);
    if (emailError) {
      setEmailError("");
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e?.preventDefault();
    e?.stopPropagation();
    console.log("1. handleSubmit dipanggil");

    setError("");
    setLoading(true);
    const toastId = toast.loading("Sedang masuk...");
    console.log("2. toast loading muncul, toastId:", toastId);

    try {
      console.log("3. sebelum login");
      await login({ email, password });
      console.log("4. login berhasil");

      toast.success("Sign In Success", {
        id: toastId,
      });

      await new Promise((resolve) => setTimeout(resolve, 1000)); // ← tunggu 1 detik

      navigate("/dashboard/outbound");
    } catch (error: any) {
      console.log("5. login gagal, error:", error);

      toast.error("Can't Sign In", {
        id: toastId,
        description: error?.message ?? "Email atau password salah.",
      });
    } finally {
      console.log("6. finally");
      setLoading(false);
    }
  };

  return (
    <div className="grid min-h-svh lg:grid-cols-2">
      <div className="relative hidden lg:block bg-primary-main overflow-hidden">
        {/* Content */}
        <div className="relative h-full flex flex-col justify-center items-center p-12 text-white">
          <div className="w-107">
            <div className="flex flex-row gap-3 items-center mb-18">
              <div className="bg-white rounded-xl w-14 h-14 flex items-center justify-center">
                <AppsIcon />
              </div>
              <h1 className="text-ml font-light tracking-wide">WMSpaceIO</h1>
            </div>
            <div>
              {/* Hero Content */}
              <div className="flex flex-col gap-2">
                <p className="text-sm font-normal tracking-widest uppercase opacity-70">
                  WMS DASHBOARD
                </p>

                <h2 className="text-large font-bold leading-tight">
                  Manage your{" "}
                  <span className="text-primary-surface">order</span> with
                  <br />
                  clarity.
                </h2>
              </div>
            </div>
            <p className="text-sm leading-5.5 text-primary-surface font-semibold pt-6">
              Track orders, manage orders, and streamline operations — all in
              one place.
            </p>
          </div>
        </div>
      </div>

      <div className="relative flex flex-col gap-4 p-6 md:p-10">
        <Toaster
          position="top-right"
          style={{
            top: 40,
            right: 180,
            left: "auto",
            width: "fit-content",
          }}
        />
        <div className="flex flex-1 items-center justify-center">
          <div className="w-123">
            <div className="flex flex-col gap-0 pb-12">
              <h1 className="text-large font-bold text-neutral-100">
                Welcome Back
              </h1>
              <h2 className="text-sm font-light text-neutral-100">
                Sign in to your account to continue
              </h2>
            </div>
            <form
              className="flex flex-col gap-8"
              onSubmit={handleSubmit}
              noValidate
            >
              <div className="flex flex-col gap-1">
                <Label htmlFor="email">Email Address</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="email@example.com"
                  value={email}
                  onChange={handleEmailChange}
                  onBlur={handleEmailBlur}
                  className={emailError ? "border-danger-main" : ""}
                />
                {emailError && (
                  <p className="text-danger-main text-xss mt-1 font-jakarta">
                    {emailError}
                  </p>
                )}
              </div>
              <div className="flex flex-col gap-1">
                <Label htmlFor="password">Password</Label>
                <Input
                  id="password"
                  type="password"
                  placeholder="Your password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                />
                {/* Remember me + Forgot password */}
                <div className="flex items-center justify-between pt-4">
                  <label className="flex items-center gap-1 cursor-pointer select-none">
                    <input
                      type="checkbox"
                      checked={rememberMe}
                      onChange={(e) => setRememberMe(e.target.checked)}
                      className="accent-primary-main size-4 cursor-pointer"
                    />
                    <span className="text-sm text-neutral-100">
                      Remember me
                    </span>
                  </label>
                  <a
                    href="/forgot-password"
                    className="text-sm font-medium text-neutral-100 hover:text-primary-main transition-colors"
                  >
                    Forgot password?
                  </a>
                </div>
              </div>
              {error && <p className="text-danger-main text-sm">{error}</p>}
              <Button
                type="submit"
                disabled={loading}
                size="lg"
                className="w-full"
                onClick={handleSubmit}
              >
                {loading ? "Loading..." : "Sign in to Dashboard"}
              </Button>

              <FieldSeparator>or</FieldSeparator>

              <Button
                type="button"
                disabled={loading}
                variant="outline"
                size="lg"
                className="w-full"
              >
                <GoogleIcon />
                Continue with Google
              </Button>

              {/* tambahan */}
              <p className="text-center text-sm font-normal text-neutral-80">
                Don't have an account?{" "}
                <a
                  href="/register"
                  className="text-info-main hover:text-primary-pressed font-medium transition-colors"
                >
                  Sign up free
                </a>
              </p>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
