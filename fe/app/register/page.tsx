"use client";

import Image from "next/image";
import Link from "next/link";
import MainPage from "@/app/components/ui/MainPage";
import FormRegister from "@/app/components/register/FormRegister";

export default function RegisterPage() {
  return (
    <MainPage>
      <div className="relative flex min-h-screen flex-col items-center justify-center bg-page px-4 py-10 text-primary-blue sm:px-8 sm:py-14">
        <header className="mb-7 flex items-center justify-center gap-3 text-center sm:gap-5">
          <Image
            src="/images/book-icon.png"
            width={100}
            height={100}
            alt=""
            className="h-20 w-20 object-contain sm:h-24 sm:w-24"
            priority
          />
          <div>
            <Link
              href="/welcome"
              className="inline-block rounded-lg text-2xl leading-tight font-black tracking-wide focus-visible:outline-2 focus-visible:outline-offset-4 sm:text-3xl"
            >
              VOCAB
              <br />
              CHALLENGE
            </Link>
            <p className="mt-2 text-xs text-[#737784] sm:text-sm">
              A little practice. A new word. Every day.
            </p>
          </div>
        </header>

        <section
          aria-labelledby="register-title"
          className="w-full max-w-xl rounded-3xl border border-white bg-white p-6 shadow-[0_16px_60px_#1b30590a] sm:p-10"
        >
          <div className="text-center">
            <span className="mb-3 inline-block rounded-full bg-amber-50 px-3 py-1 text-xs font-semibold tracking-wide text-amber-700">
              YOUR NEXT CHAPTER
            </span>
            <h1
              id="register-title"
              className="text-3xl font-bold tracking-tight sm:text-4xl"
            >
              Create an account
            </h1>
            <p className="mt-3 text-sm leading-relaxed text-[#737784]">
              Start your vocab journey and
              <br />
              get a little better every day.
            </p>
          </div>

          <FormRegister />
          <p className="mt-6 text-center text-sm text-[#737784]">
            Already have an account?{" "}
            <Link
              href="/login"
              className="ml-1 font-semibold text-[#0861db]! hover:underline focus-visible:outline-2 focus-visible:outline-offset-4"
            >
              Log In
            </Link>
          </p>
        </section>
      </div>
    </MainPage>
  );
}
