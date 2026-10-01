import { VocabularyDemo } from "./components/VocabularyDemo";
import Link from "next/link";

export default function Home() {
  return (
    <main className="min-h-screen px-5 py-10 sm:px-8">
      <section className="mx-auto flex w-full max-w-3xl flex-col gap-6">
        <div className="space-y-2">
          <p className="text-sm font-semibold tracking-wide text-teal-700 uppercase">
            Vocabulary FE
          </p>
          <h1 className="text-3xl font-bold text-slate-950 sm:text-4xl">
            Frontend stack installed and connected.
          </h1>
          <p className="max-w-2xl text-base text-slate-600">
            This starter is ready for building vocabulary features with a modern
            Next.js app shell.
          </p>
          <Link
            href="/vocabulary"
            className="inline-flex text-sm font-semibold text-blue-700 hover:text-blue-900"
          >
            Open vocabulary table
          </Link>
        </div>
        <VocabularyDemo />
      </section>
    </main>
  );
}
