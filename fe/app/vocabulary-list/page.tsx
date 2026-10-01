import Link from "next/link";
import { VocabularyTable } from "./VocabularyTable";

export default function VocabularyPage() {
  return (
    <main className="min-h-screen px-5 py-8 sm:px-8">
      <section className="mx-auto flex w-full max-w-6xl flex-col gap-6">
        <div className="flex flex-col justify-between gap-4 md:flex-row md:items-end">
          <div className="space-y-2">
            <p className="text-sm font-semibold tracking-wide text-teal-700 uppercase">
              Vocabulary
            </p>
            <h1 className="text-3xl font-bold text-slate-950 sm:text-4xl">
              Vocabulary Table
            </h1>
            <p className="max-w-2xl text-base text-slate-600">
              Review words, meanings, levels, examples, and latest practice
              dates in one place.
            </p>
          </div>
          <Link
            href="/"
            className="text-sm font-semibold text-blue-700 hover:text-blue-900"
          >
            Back to home
          </Link>
        </div>

        <VocabularyTable />
      </section>
    </main>
  );
}
