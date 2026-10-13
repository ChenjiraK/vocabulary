import api from "./api";

type GoogleLoginResponse = {
  url: string;
};

export type AuthUser = {
  id: number;
  display_name: string;
  avatar_url?: string;
  user_type: string;
};

export type AuthLoginResponse = {
  message: string;
  token: string;
  user: AuthUser;
};

export async function loginWithGoogle() {
  const { data } = await api.get<GoogleLoginResponse>("/auth/google/login");
  return data;
}

export async function loginWithFacebook() {
  const { data } = await api.get<GoogleLoginResponse>("/auth/facebook/login");
  return data;
}

export async function loginWithPassword(payload: {
  email: string;
  password: string;
}) {
  const { data } = await api.post<AuthLoginResponse>("/auth/login", {
    email: payload.email,
    password: payload.password,
  });
  return data;
}
