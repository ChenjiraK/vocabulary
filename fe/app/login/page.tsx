"use client";

import { type FormEvent, useState } from "react";
import Image from "next/image";
import { useRouter } from "next/navigation";
import {
  loginWithFacebook,
  loginWithGoogle,
  loginWithPassword,
} from "@/app/services/auth.service";
import { useAuthStore } from "@/app/stores/auth.store";
import MainPage from "@/app/components/ui/MainPage";
import FormEmailPassword from "@/app/components/ui/form/FormEmailPassword";
import GoogleButton from "@/app/components/ui/form/GoogleButton";
import FacebookButton from "@/app/components/ui/form/FacebookButton";

export default function LoginPage() {
  const router = useRouter();
  const setLoginData = useAuthStore((state) => state.setLoginData);
  const [showPassword, setShowPassword] = useState(false);
  const [pending, setPending] = useState<
    "password" | "google" | "facebook" | null
  >(null);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (pending) return;
    const form = new FormData(event.currentTarget);
    setError("");
    setNotice("");
    setPending("password");
    try {
      const data = await loginWithPassword({
        email: String(form.get("username")).trim(),
        password: String(form.get("password")),
      });
      setLoginData(data);
      router.push("/game");
    } catch {
      setError("Unable to log in. Please check your details and try again.");
    } finally {
      setPending(null);
    }
  }

  async function handleSocialLogin(provider: "google" | "facebook") {
    if (pending) return;
    setError("");
    setNotice("");
    setPending(provider);
    try {
      const data = await (provider === "google"
        ? loginWithGoogle()
        : loginWithFacebook());
      const url = new URL(data.url);
      if (url.protocol !== "https:" && url.protocol !== "http:")
        throw new Error("Invalid login URL");
      window.location.assign(url.href);
    } catch {
      setError(
        `Unable to connect to ${provider === "google" ? "Google" : "Facebook"}. Please try again.`
      );
      setPending(null);
    }
  }

  return (
    <MainPage>
      <div className="bg-page flex min-h-screen items-center justify-center px-8 py-10">
        <main>
          <section className="flex">
            <div>
              <Image
                src="/images/book-icon.png"
                alt="Correct answer"
                width={200}
                height={200}
                unoptimized
              />
            </div>
            <div className="flex w-full max-w-[320px] flex-col items-center text-center">
              <h1 className="text-[48px] leading-[1.1] font-black tracking-wide text-primary-blue">
                VOCAB
                <br />
                CHALLENGE
              </h1>
              <p className="text-primary-blue pt-2">
                Fun Vocabulary Games! <br /> Practice Every Day, Get Better
                Every Day!
              </p>
            </div>
            <div>
              <Image
                src="/images/star-icon.png"
                alt="Correct answer"
                width={150}
                height={150}
                unoptimized
              />
            </div>
          </section>
          <section className="bg-white rounded-3xl shadow-xl p-10 max-w-xl">
            <h2 id="login-title" className="font-bold text-4xl text-center">
              Welcome !
            </h2>
            <p className="text-gray-500 text-center">
              Log in to continue your vocab journey
              <br /> and keep challenging yourself.
            </p>
            <div>
              <GoogleButton />
              <FacebookButton />
              <div className="my-6 flex items-center gap-4">
                <div className="flex-1 border-t border-gray-300" />
                <span className="text-sm text-gray-500">or</span>
                <div className="flex-1 border-t border-gray-300" />
              </div>
              <FormEmailPassword />
            </div>
          </section>
        </main>
      </div>
    </MainPage>
  );
}
