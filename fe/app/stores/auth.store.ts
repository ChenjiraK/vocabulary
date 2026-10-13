"use client";

import { create } from "zustand";
import { persist } from "zustand/middleware"; //บันทึก state ลง storage
import type { AuthLoginResponse, AuthUser } from "@/app/services/auth.service";

type AuthState = {
  message: string | null;
  token: string | null;
  user: AuthUser | null;
  isAuthenticated: boolean;
  setLoginData: (data: AuthLoginResponse) => void;
  clearAuth: () => void;
};

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      message: null,
      token: null,
      user: null,
      isAuthenticated: false,
      setLoginData: (data) =>
        set({
          message: data.message,
          token: data.token,
          user: data.user,
          isAuthenticated: Boolean(data.token),
        }),
      clearAuth: () =>
        set({
          message: null,
          token: null,
          user: null,
          isAuthenticated: false,
        }),
    }),
    {
      name: "auth-storage",
    }
  )
);
