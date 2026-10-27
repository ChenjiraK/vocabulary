"use client";

import { useState } from "react";
import { FormProvider, type SubmitHandler, useForm } from "react-hook-form";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faEnvelope,
  faEye,
  faEyeSlash,
  faLock,
  faUser,
} from "@fortawesome/free-solid-svg-icons";
import ProfilePicture from "./ProfilePicture";

export type RegisterFormValues = {
  profilePicture: File | null;
  displayName: string;
  email: string;
  password: string;
  confirmPassword: string;
};

const inputClass =
  "min-w-0 w-full rounded-xl border-0 bg-transparent px-3 py-3 text-base text-primary-blue outline-none placeholder:text-gray-400";
const wrapperClass =
  "flex min-h-12 items-center rounded-xl border-2 border-[#dadddf] bg-white transition-colors focus-within:border-[#5678af] focus-within:ring-3 focus-within:ring-[#5678af18] motion-reduce:transition-none";
const labelClass = "mb-2 block text-sm font-medium text-[#626775]";

export default function FormRegister() {
  const [notice, setNotice] = useState("");
  const [visible, setVisible] = useState({
    password: false,
    confirmPassword: false,
  });
  const form = useForm<RegisterFormValues>({
    mode: "onChange",
    defaultValues: {
      profilePicture: null,
      displayName: "",
      email: "",
      password: "",
      confirmPassword: "",
    },
  });
  const {
    register,
    handleSubmit,
    getValues,
    formState: { errors, isSubmitting },
  } = form;

  const onSubmit: SubmitHandler<RegisterFormValues> = () => {
    setNotice("");
    // TODO: Connect the registration API using the validated form values.
  };

  return (
    <FormProvider {...form}>
      <form
        noValidate
        onSubmit={handleSubmit(onSubmit)}
        onChange={() => {
          setNotice("");
        }}
        className="mt-7 space-y-5 [&_button]:focus-visible:outline-2 [&_button]:focus-visible:outline-offset-4 [&_button]:focus-visible:outline-[#5678af]"
      >
        <ProfilePicture />

        <div>
          <label htmlFor="displayName" className={labelClass}>
            Display name
          </label>
          <div className={wrapperClass}>
            <FontAwesomeIcon
              icon={faUser}
              className="ml-4 w-4 shrink-0 text-[#7d818b]"
            />
            <input
              id="displayName"
              {...register("displayName", {
                required: "Please enter your display name.",
                setValueAs: (value: string) => value.trim(),
                maxLength: {
                  value: 80,
                  message: "Display name must be 80 characters or fewer.",
                },
              })}
              autoComplete="nickname"
              placeholder="What should we call you?"
              required
              maxLength={80}
              aria-invalid={!!errors.displayName}
              aria-describedby="displayName-error"
              className={inputClass}
            />
          </div>
          <p
            id="displayName-error"
            role="alert"
            className="mt-2 text-xs text-red-600 empty:hidden"
          >
            {errors.displayName?.message}
          </p>
        </div>
        <div>
          <label htmlFor="email" className={labelClass}>
            Email
          </label>
          <div className={wrapperClass}>
            <FontAwesomeIcon
              icon={faEnvelope}
              className="ml-4 w-4 shrink-0 text-[#7d818b]"
            />
            <input
              id="email"
              {...register("email", {
                required: "Please enter your email.",
                setValueAs: (value: string) => value.trim(),
                pattern: {
                  value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
                  message: "Please enter a valid email address.",
                },
              })}
              type="email"
              autoComplete="email"
              autoCapitalize="none"
              spellCheck={false}
              placeholder="Enter your email"
              required
              aria-invalid={!!errors.email}
              aria-describedby="email-error"
              className={inputClass}
            />
          </div>
          <p
            id="email-error"
            role="alert"
            className="mt-2 text-xs text-red-600 empty:hidden"
          >
            {errors.email?.message}
          </p>
        </div>
        {(["password", "confirmPassword"] as const).map((name) => (
          <div key={name}>
            <label htmlFor={name} className={labelClass}>
              {name === "password" ? "Password" : "Confirm password"}
            </label>
            <div className={wrapperClass}>
              <FontAwesomeIcon
                icon={faLock}
                className="ml-4 w-4 shrink-0 text-[#7d818b]"
              />
              <input
                id={name}
                {...register(name, {
                  required:
                    name === "password"
                      ? "Please enter your password."
                      : "Please confirm your password.",
                  minLength: {
                    value: 8,
                    message: "Use at least 8 characters.",
                  },
                  ...(name === "password"
                    ? { deps: ["confirmPassword"] }
                    : {
                        validate: (value: string) =>
                          value === getValues("password") ||
                          "Passwords do not match. Please try again.",
                      }),
                })}
                type={visible[name] ? "text" : "password"}
                autoComplete="new-password"
                placeholder={
                  name === "password"
                    ? "Create a password"
                    : "Re-enter your password"
                }
                required
                minLength={8}
                aria-describedby={
                  name === "password"
                    ? "password-hint password-error"
                    : "confirmPassword-error"
                }
                aria-invalid={!!errors[name]}
                className={inputClass}
              />
              <button
                type="button"
                aria-label={`${visible[name] ? "Hide" : "Show"} ${name === "password" ? "password" : "confirm password"}`}
                aria-pressed={visible[name]}
                aria-controls={name}
                className="grid min-h-11 w-11 shrink-0 place-items-center rounded-xl text-[#7d818b] hover:text-primary-blue"
                onClick={() =>
                  setVisible((state) => ({
                    ...state,
                    [name]: !state[name],
                  }))
                }
              >
                <FontAwesomeIcon icon={visible[name] ? faEye : faEyeSlash} />
              </button>
            </div>
            {name === "password" && (
              <p id="password-hint" className="mt-2 text-xs text-[#737784]">
                Use at least 8 characters.
              </p>
            )}
            <p
              id={`${name}-error`}
              role="alert"
              className="mt-2 text-xs text-red-600 empty:hidden"
            >
              {errors[name]?.message}
            </p>
          </div>
        ))}
        <button
          type="submit"
          disabled={isSubmitting}
          className="min-h-12 mt-3 w-full rounded-xl bg-primary-blue px-4 py-3 font-bold text-white shadow-sm transition-colors hover:bg-[#294777] active:bg-[#102a53] motion-reduce:transition-none"
        >
          Create account
        </button>
        <p
          role="status"
          className="text-center text-sm text-[#737784] empty:hidden"
        >
          {notice}
        </p>
      </form>
    </FormProvider>
  );
}
