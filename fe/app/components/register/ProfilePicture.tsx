"use client";

import { type ChangeEvent, useEffect, useRef, useState } from "react";
import { useController, useFormContext } from "react-hook-form";
import Image from "next/image";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCamera, faUser } from "@fortawesome/free-solid-svg-icons";
import type { RegisterFormValues } from "./FormRegister";

export default function ProfilePicture() {
  const [preview, setPreview] = useState("");
  const fileInput = useRef<HTMLInputElement>(null);
  const unreadablePicture = useRef<File | null>(null);
  const {
    control,
    getValues,
    resetField,
    setError,
    formState: { errors },
  } = useFormContext<RegisterFormValues>();
  const { field: profilePicture } = useController({
    name: "profilePicture",
    control,
    rules: {
      validate: {
        format: (file) =>
          !file ||
          (["image/jpeg", "image/png", "image/webp"].includes(file.type) &&
            file.size <= 5 * 1024 * 1024) ||
          "Please choose a JPG, PNG or WebP image under 5 MB.",
        readable: (file) =>
          !file ||
          file !== unreadablePicture.current ||
          "This image could not be opened. Please choose another image.",
      },
    },
  });

  useEffect(
    () => () => {
      if (preview) URL.revokeObjectURL(preview);
    },
    [preview],
  );

  function handlePicture(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    if (!file) return;
    unreadablePicture.current = null;
    profilePicture.onChange(file);
    if (
      !["image/jpeg", "image/png", "image/webp"].includes(file.type) ||
      file.size > 5 * 1024 * 1024
    ) {
      setPreview("");
      return;
    }
    setPreview(URL.createObjectURL(file));
  }

  return (
    <div className="flex flex-col items-center rounded-2xl bg-[#f8faff] px-4 py-5">
      <label
        htmlFor="profilePicture"
        className="group relative cursor-pointer rounded-full focus-within:ring-3 focus-within:ring-[#5678af] focus-within:ring-offset-4"
      >
        <span className="flex h-22 w-22 items-center justify-center overflow-hidden rounded-full border-4 border-white bg-[#e9eff9] shadow-sm">
          {preview ? (
            <Image
              src={preview}
              width={88}
              height={88}
              unoptimized
              alt="Profile picture preview"
              className="h-full w-full object-cover"
              onError={() => {
                setPreview("");
                unreadablePicture.current = getValues("profilePicture");
                setError("profilePicture", {
                  type: "validate",
                  message:
                    "This image could not be opened. Please choose another image.",
                });
                if (fileInput.current) fileInput.current.value = "";
              }}
            />
          ) : (
            <FontAwesomeIcon
              icon={faUser}
              className="text-3xl text-[#8b9cbb]"
            />
          )}
        </span>
        <span className="absolute right-0 bottom-0 grid h-8 w-8 place-items-center rounded-full border-2 border-white bg-primary-blue text-xs text-white transition-colors group-hover:bg-[#294777]">
          <FontAwesomeIcon icon={faCamera} />
        </span>
        <input
          ref={(element) => {
            profilePicture.ref(element);
            fileInput.current = element;
          }}
          id="profilePicture"
          name={profilePicture.name}
          type="file"
          accept="image/jpeg,image/png,image/webp"
          aria-label="Upload profile picture"
          aria-describedby="picture-hint picture-error"
          aria-invalid={!!errors.profilePicture}
          className="sr-only"
          onChange={handlePicture}
          onBlur={profilePicture.onBlur}
        />
      </label>
      <p className="mt-3 text-sm font-semibold">
        Profile picture{" "}
        <span className="font-normal text-gray-400">(optional)</span>
      </p>
      <p id="picture-hint" className="mt-1 text-xs text-[#737784]">
        Click the camera to upload · JPG, PNG, WebP · Max 5 MB
      </p>
      {preview && (
        <button
          type="button"
          className="mt-2 text-xs text-[#0861db] hover:underline"
          onClick={() => {
            setPreview("");
            unreadablePicture.current = null;
            resetField("profilePicture");
            if (fileInput.current) fileInput.current.value = "";
          }}
        >
          Remove picture
        </button>
      )}
      <p
        id="picture-error"
        role="alert"
        className="text-center text-xs text-red-600 empty:hidden"
      >
        {errors.profilePicture?.message}
      </p>
    </div>
  );
}
