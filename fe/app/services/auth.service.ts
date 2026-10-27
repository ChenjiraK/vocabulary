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
  correct_answers: number;
  score: number;
  message: string;
  token: string;
  user: AuthUser;
};

export async function loginWithGoogle() {
  try {
    const { data } = await api.get<GoogleLoginResponse>("/auth/google/login");
    return data;
  } catch (error) {
    console.error("Error logging in with Google:", error);
    throw error;
  }
}

export async function loginWithFacebook() {
  try {
    const { data } = await api.get<GoogleLoginResponse>("/auth/facebook/login");
    return data;
  } catch (error) {
    console.error("Error logging in with Facebook:", error);
    throw error;
  }
}

export async function loginWithPassword(payload: {
  email: string;
  password: string;
}) {
  return new Promise<AuthLoginResponse>(async (resolve, reject) => {
    try {
      const { data } = await api.post<AuthLoginResponse>("/auth/login", {
        email: payload.email,
        password: payload.password,
      });
      resolve(data);
    } catch (error) {
      console.error("Error logging in with password:", error);
      reject(error);
    }
  });
}

export async function getCurrentUser() {
  const { data } = await api.get<AuthLoginResponse>("/auth/me");
  return data;
}

export async function register(payload: {
  avatar_url?: string;
  display_name: string;
  email: string;
  password: string;
}) {
  const { data } = await api.post<any>("/auth/register", { ...payload });
  return data;
}
