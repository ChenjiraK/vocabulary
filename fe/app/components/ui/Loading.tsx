"use client";

import { useId } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faBookOpen, faStar } from "@fortawesome/free-solid-svg-icons";

export type LoadingProps = {
  title?: string;
  message?: string;
};

export default function Loading({
  title = "Just a moment…",
  message = "Getting everything ready for your vocab journey.",
}: LoadingProps) {
  const id = useId();

  return (
    <div
      role="status"
      aria-live="polite"
      aria-atomic="true"
      aria-labelledby={`${id}-title`}
      aria-describedby={`${id}-description`}
      className="relative isolate w-full max-w-sm overflow-hidden rounded-[28px] border border-white bg-[#fffdf8] px-6 py-9 text-center text-primary-blue shadow-[0_16px_60px_#1b30590a] sm:px-9 sm:py-10"
    >
      <div
        aria-hidden="true"
        className="pointer-events-none absolute -top-20 -right-20 -z-10 h-48 w-48 rounded-full bg-[#fff1cc]/60"
      />
      <div
        aria-hidden="true"
        className="pointer-events-none absolute -bottom-24 -left-20 -z-10 h-44 w-44 rounded-full bg-[#e7f0fc]/60"
      />

      <p
        aria-hidden="true"
        className="text-[10px] font-bold tracking-[0.2em] text-[#8190a7]"
      >
        VOCAB CHALLENGE
      </p>

      <div
        aria-hidden="true"
        className="relative mx-auto mt-7 mb-7 grid h-36 w-44 place-items-center"
      >
        <span className="absolute h-32 w-32 rounded-full border-[3px] border-[#e6edf7]" />
        <span className="absolute h-32 w-32 animate-spin rounded-full border-[3px] border-transparent border-t-primary-blue border-r-[#f6c647] [animation-duration:1.5s] motion-reduce:animate-none" />
        <span className="grid h-24 w-24 place-items-center rounded-full bg-[#fff0c7] shadow-[inset_0_0_0_5px_#fff8e6]">
          <FontAwesomeIcon
            icon={faBookOpen}
            className="text-4xl text-primary-blue"
          />
        </span>
        <FontAwesomeIcon
          icon={faStar}
          className="absolute top-0 right-2 rotate-12 text-lg text-[#f5c349]"
        />
        <span className="absolute bottom-2 left-2 flex h-7 w-7 -rotate-12 items-center justify-center rounded-lg bg-white text-sm font-bold text-[#6a89b3] shadow-sm">
          A
        </span>
        <span className="absolute top-4 left-2 h-2 w-2 rounded-full bg-[#b5ccec]" />
        <span className="absolute right-1 bottom-6 h-1.5 w-1.5 rounded-full bg-[#efd999]" />
      </div>

      <h2
        id={`${id}-title`}
        className="text-2xl leading-tight font-bold tracking-tight"
      >
        {title}
      </h2>
      <p
        id={`${id}-description`}
        className="mx-auto mt-3 max-w-64 text-sm leading-relaxed text-[#737784]"
      >
        {message}
      </p>

      <div aria-hidden="true" className="mt-6 flex justify-center gap-2">
        <span className="h-2 w-2 animate-bounce rounded-full bg-primary-blue [animation-duration:1.2s] motion-reduce:animate-none" />
        <span className="h-2 w-2 animate-bounce rounded-full bg-[#7895bf] [animation-delay:150ms] [animation-duration:1.2s] motion-reduce:animate-none" />
        <span className="h-2 w-2 animate-bounce rounded-full bg-[#f5c349] [animation-delay:300ms] [animation-duration:1.2s] motion-reduce:animate-none" />
      </div>
    </div>
  );
}
