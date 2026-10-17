"use client";

import { type FormEvent, useState } from "react";
import { faEye, faEyeSlash } from "@fortawesome/free-regular-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { IconButton, InputAdornment, TextField } from "@mui/material";
import { useRouter } from "next/navigation";
import { loginWithPassword } from "@/app/services/auth.service";
import { useAuthStore } from "@/app/stores/auth.store";

interface FormEmailPasswordProps {
  onBack?: () => void;
}

export default function FormEmailPassword({ onBack }: FormEmailPasswordProps) {
  const [email, setEmail] = useState("test@example.com");
  const [password, setPassword] = useState("password");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const router = useRouter();
  const setLoginData = useAuthStore((state) => state.setLoginData);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setIsLoading(true);

    try {
      const payload = {
        email,
        password,
      };
      const data = await loginWithPassword(payload);
      setLoginData(data);
      router.push("/game");
    } catch (err) {
      console.error(err);
      setError("Login failed. Please check your email and password.");
    } finally {
      setIsLoading(false);
    }
  }
  return (
    <div>
      <h2 className="text-lg font-semibold mb-4">
        Login with Email and Password
      </h2>
      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <TextField
          label="Email"
          type="email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          required
          fullWidth
          size="small"
          sx={{
            my: 1.5,
          }}
        />
        <TextField
          label="Password"
          type={showPassword ? "text" : "password"}
          value={password}
          onChange={(event) => setPassword(event.target.value)}
          required
          fullWidth
          size="small"
          slotProps={{
            input: {
              endAdornment: (
                <InputAdornment position="end">
                  <IconButton
                    aria-label={
                      showPassword ? "Hide password" : "Show password"
                    }
                    edge="end"
                    onClick={() => setShowPassword((value) => !value)}
                  >
                    <FontAwesomeIcon
                      icon={showPassword ? faEyeSlash : faEye}
                      className="text-base"
                    />
                  </IconButton>
                </InputAdornment>
              ),
            },
          }}
        />
        {error && <p className="text-sm text-red-600">{error}</p>}
        <div className="flex gap-5 mt-5">
          <button
            type="submit"
            disabled={isLoading}
            className="flex h-10 w-full items-center justify-center gap-2 rounded-[40px] border border-[#b9c0c9] bg-primary-blue text-sm font-bold text-white shadow-sm"
          >
            Login
          </button>
          <button
            type="button"
            onClick={onBack}
            className="flex h-10 w-full items-center justify-center gap-2 rounded-[40px] border border-[#b9c0c9] bg-white text-sm font-bold text-slate-950 shadow-sm"
          >
            Back
          </button>
        </div>
      </form>
    </div>
  );
}
