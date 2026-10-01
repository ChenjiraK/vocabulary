"use client";

import { type ChangeEvent, type FormEvent, useState } from "react";
import MainPage from "../components/MainPage";
import CircleCountdown from "../components/CircleCountdown";

export type GameWord = {
  id: number;
  vocab: string;
  meaning: string;
  synonyms: string;
  status: "" | "success" | "failed";
  thai: string;
  type_of_speech: string;
};

type WordInputProps = {
  word: GameWord;
  degree: number;
  onSubmit: (inputWord: string) => void;
};

export default function WordInput({ word, degree, onSubmit }: WordInputProps) {
  const [inputWord, setInputWord] = useState("");
  const [isShowHint, setIsShowHint] = useState(false);
  const answer = word.vocab;
  const answerCharacters = Array.from(answer);

  function handleInputChange(event: ChangeEvent<HTMLInputElement>) {
    const nextValue = event.target.value
      .toUpperCase()
      .replace(/[^A-Z]/g, "")
      .slice(0, answerCharacters.length);

    setInputWord(nextValue);

    if (nextValue.length === answerCharacters.length) {
      onSubmit(nextValue);
    }
  }

  function classInputWord() {
    if (inputWord.length === answerCharacters.length) {
      if (inputWord === answer.toUpperCase()) {
        return "border-[#4cb131]";
      } else {
        return "border-[#f02b5d]";
      }
    }
    return "border-[#aeb4bd]";
  }

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    onSubmit(inputWord);
  }

  return (
    <MainPage>
      <div className="bg-page flex min-h-180 flex-col rounded-b-md">
        <div className="grid grid-cols-[44px_1fr_52px] items-center px-4 pt-8">
          <button
            type="button"
            aria-label="back"
            className="flex h-9 w-9 items-center justify-center text-[#111827]"
          >
            <span className="block h-3 w-3 rotate-45 border-b-2 border-l-2 border-current" />
          </button>

          <p className="text-primary-blue text-center text-3xl font-semibold">
            Level {degree}
          </p>
          <CircleCountdown />
        </div>

        <section className="flex flex-1 flex-col items-center px-5 pt-16 text-center">
          <p className="text-[17px] leading-6 font-medium text-slate-950">
            <span className="text-primary-blue font-bold">
              ({word.type_of_speech})
            </span>{" "}
            {word.meaning}
          </p>

          <form className="relative mt-16 w-full px-5" onSubmit={handleSubmit}>
            <label
              htmlFor="input-word"
              className="flex flex-wrap justify-center gap-2"
            >
              {answerCharacters.map((character, index) => (
                <div
                  key={`${character}-${index}`}
                  className={`flex aspect-square min-h-12 w-12 max-w-18 min-w-12 items-center justify-center rounded-md border-2 bg-white text-2xl font-bold text-slate-950 shadow-[inset_0_1px_0_rgba(255,255,255,0.65)] sm:h-18 sm:w-18 ${classInputWord()}`}
                >
                  {inputWord[index] ?? ""}
                </div>
              ))}
            </label>
            <input
              type="text"
              id="input-word"
              value={inputWord}
              maxLength={answerCharacters.length}
              autoComplete="off"
              aria-label="answer"
              className="absolute inset-0 h-full w-full cursor-text opacity-0"
              onChange={handleInputChange}
            />
          </form>
          {!isShowHint && (
            <button
              type="button"
              className="mt-20 flex h-10 w-36.25 items-center justify-center gap-2 rounded-md border border-[#b9c0c9] bg-white text-sm font-bold text-slate-950 shadow-sm"
              onClick={() => setIsShowHint(true)}
            >
              <span aria-hidden="true" className="text-base">
                💡
              </span>
              HINT
            </button>
          )}

          {isShowHint && (
            <div className="text-primary-blue mt-20 items-center justify-center rounded-2xl border-2 border-amber-300 bg-[#fdfae1] p-4 text-center font-medium">
              <div>
                start with{" "}
                <span className="font-bold text-emerald-700">
                  {word.vocab[0].toUpperCase()}
                </span>{" "}
                and end with{" "}
                <span className="font-bold text-emerald-700">
                  {word.vocab[word.vocab.length - 1].toUpperCase()}
                </span>
              </div>
              <div className="mt-3 py-3 text-left font-bold">💡 meaning </div>
              {word.thai ? word.thai : word.synonyms}
            </div>
          )}
        </section>
      </div>
    </MainPage>
  );
}
