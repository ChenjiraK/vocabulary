"use client";

import { type FormEvent, useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faEye, faEyeSlash, faLock } from "@fortawesome/free-solid-svg-icons";
import { faUser } from "@fortawesome/free-regular-svg-icons";
import { useRouter } from "next/navigation";
import { loginWithPassword } from "@/app/services/auth.service";
import { useAuthStore } from "@/app/stores/auth.store";

interface FormEmailPasswordProps {
  onBack?: () => void;
  setNotice?: (message: string) => void;
  pending?: boolean;
}

export default function FormEmailPassword({
  pending = false,
  setNotice,
}: FormEmailPasswordProps) {
  const [email, setEmail] = useState("test@example.com");
  const [password, setPassword] = useState("password");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const router = useRouter();
  const setLoginData = useAuthStore((state) => state.setLoginData);

  const isPending = pending || isLoading;

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (isPending) return;
    setError("");
    setNotice?.("");
    setIsLoading(true);

    try {
      const payload = {
        email: email.trim(),
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
  //NOTE: className="h-[clamp(20px,3vw,31px)]! w-[clamp(18px,2.7vw,28px)]!"
  //ความสูงจะปรับตาม 3vw (3% ของความกว้าง viewport) แต่ต้องไม่ต่ำกว่า 20px และไม่เกิน 31px บังคับใช้ !important เพื่อเพิ่มความสำคัญของ CSS
  return (
    <div>
      <form
        className="flex flex-col [&_button:disabled]:cursor-wait [&_button:disabled]:opacity-60 [&_button:focus-visible]:outline-3 [&_button:focus-visible]:outline-offset-4 [&_button:focus-visible]:outline-[#668dcc] motion-reduce:[&_button]:transition-none"
        onSubmit={handleSubmit}
        aria-busy={isPending}
      >
        <div>
          <label
            htmlFor="email"
            className="mb-2.25 block leading-[1.2] font-medium text-[#626775]"
          >
            Email
          </label>
          <div className="flex h-11 items-center rounded-xl border-2 border-[#dadddf] bg-white transition-[border-color,box-shadow] duration-180 focus-within:border-[#5678af] focus-within:shadow-[0_0_0_3px_#5678af18] motion-reduce:transition-none">
            <FontAwesomeIcon
              icon={faUser}
              className="ml-3.5 h-5 w-5 shrink-0 text-[#7d818b]"
            />
            <input
              id="email"
              name="email"
              type="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              autoComplete="username"
              autoCapitalize="none"
              spellCheck={false}
              className="min-w-0 w-full self-stretch rounded-[inherit] border-0 bg-transparent px-3 py-3 text-[#132b50] outline-none placeholder:text-[#969ba5] placeholder:opacity-100"
              placeholder="Enter your email"
              required
              disabled={isPending}
            />
          </div>
        </div>
        <div className="mt-5">
          <label
            htmlFor="password"
            className="mb-2.25 block leading-[1.2] font-medium text-[#626775]"
          >
            Password
          </label>
          <div className="flex h-11 items-center rounded-xl border-2 border-[#dadddf] bg-white transition-[border-color,box-shadow] duration-180 focus-within:border-[#5678af] focus-within:shadow-[0_0_0_3px_#5678af18] motion-reduce:transition-none">
            <FontAwesomeIcon
              icon={faLock}
              className="ml-3.5 h-5 w-5 shrink-0 text-[#7d818b]"
            />
            <input
              id="password"
              name="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              type={showPassword ? "text" : "password"}
              autoComplete="current-password"
              className="min-w-0 w-full self-stretch rounded-[inherit] border-0 bg-transparent px-3 py-3 text-[#132b50] outline-none placeholder:text-[#969ba5] placeholder:opacity-100"
              placeholder="Enter your password"
              required
              disabled={isPending}
            />
            <button
              type="button"
              className="grid w-11 shrink-0 place-items-center self-stretch rounded-xl border-0 bg-transparent text-[20px] text-[#7d818b]"
              aria-label={showPassword ? "Hide password" : "Show password"}
              aria-pressed={showPassword}
              aria-controls="password"
              disabled={isPending}
              onClick={() => setShowPassword((visible) => !visible)}
            >
              <FontAwesomeIcon
                className="h-3.5! w-3.5!"
                icon={showPassword ? faEye : faEyeSlash}
              />
            </button>
          </div>
        </div>
        <button
          type="button"
          className="mt-3.5 self-end border-0 bg-transparent p-0 text-[12px] leading-[1.3] text-[#0861db] hover:underline"
          disabled={isPending}
          onClick={() =>
            setNotice
              ? setNotice("Password reset is not available yet.")
              : setError("Password reset is not available yet.")
          }
        >
          Forgot password?
        </button>
        {error && (
          <p className="mt-4 text-sm text-[#b52929]" role="alert">
            {error}
          </p>
        )}
        <button
          type="submit"
          className="mt-5 h-11 w-full rounded-xl border border-[#142e5a] bg-[linear-gradient(110deg,#1a3766,#19335e)] font-bold text-white shadow-[0_3px_5px_#0c21490a] transition-[background,transform] duration-180 hover:bg-[#102a53] hover:bg-none active:translate-y-px"
          disabled={isPending}
        >
          {isLoading ? "Logging in…" : "Log In"}
        </button>
      </form>
      <div className="text-gray-500 text-center mt-4">
        Don't have account?{" "}
        <span className="text-[#0861db] underline cursor-pointer">Sign Up</span>
      </div>
    </div>
  );
}
