"use client";

import { useEffect, useRef, useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faForward } from "@fortawesome/free-solid-svg-icons";
import { useAuthStore } from "@/app/stores/auth.store";
import MainPage from "../components/ui/MainPage";
import IncorrectModal from "../components/modal/IncorrectModal";
import SuccessModal from "../components/modal/SuccessModal";
import WordInput from "./WordInput";
import type { GameVocabResponse } from "@/app/type/GameType";
import { getVocab, submitVocabAnswer } from "@/app/services/game.service";

type SubmissionResult = "correct" | "incorrect" | null;

export default function GamePage() {
  const [submissionResult, setSubmissionResult] =
    useState<SubmissionResult>(null);
  const [gameResponse, setGameResponse] = useState<GameVocabResponse | null>(
    null
  );
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [round, setRound] = useState(0);
  const submitting = useRef(false);
  const level = useAuthStore.getState().correct_answer ?? 1;

  useEffect(() => {
    let cancelled = false;

    async function loadWord() {
      try {
        const word = await getVocab();
        if (!cancelled) setGameResponse(word);
      } catch {
        if (!cancelled) setError("Unable to load a word. Please try again.");
      } finally {
        if (!cancelled) setIsLoading(false);
      }
    }

    void loadWord();
    return () => {
      cancelled = true;
    };
  }, [round]);

  async function handleSubmit(inputWord: string) {
    if (
      !gameResponse ||
      submitting.current ||
      submissionResult === "correct" ||
      inputWord.length !== gameResponse.length_of_word
    )
      return;

    submitting.current = true;
    setIsSubmitting(true);
    setError(null);
    try {
      const result = await submitVocabAnswer(gameResponse.id, inputWord);
      setSubmissionResult(result.is_correct ? "correct" : "incorrect");
      setGameResponse({
        ...gameResponse,
        status: result.is_correct ? "success" : "failed",
      });
    } catch {
      setError("Unable to check your answer. Press Enter to try again.");
    } finally {
      submitting.current = false;
      setIsSubmitting(false);
    }
  }

  function handleNextWord() {
    setSubmissionResult(null);
    setGameResponse(null);
    setError(null);
    setIsLoading(true);
    // Increment the round to trigger a new word fetch
    setRound((value) => value + 1);
  }

  function goToWelcomePage() {
    window.location.href = "/welcome";
  }

  return (
    <>
      {isLoading ? (
        <MainPage>
          <p role="status" className="p-8 text-center">
            Loading word...
          </p>
        </MainPage>
      ) : gameResponse ? (
        <fieldset disabled={isSubmitting} className="m-0 min-w-0 border-0 p-0">
          <WordInput
            key={`${round}-${gameResponse.id}`}
            word={gameResponse}
            degree={level}
            onSubmit={handleSubmit}
            onNextWord={handleNextWord}
          ></WordInput>
        </fieldset>
      ) : (
        <MainPage>
          <div className="p-8 text-center">
            {!error && <p>No vocabulary available.</p>}
            <button
              type="button"
              onClick={goToWelcomePage}
              className="mt-4 underline"
            >
              Try again
            </button>
          </div>
        </MainPage>
      )}
      {error && (
        <p role="alert" className="p-4 text-center text-red-600">
          {error}
        </p>
      )}
      {isSubmitting && (
        <p role="status" className="p-4 text-center">
          Checking answer...
        </p>
      )}
      <SuccessModal
        isOpen={submissionResult === "correct"}
        onSubmit={handleNextWord}
        onClose={() => setSubmissionResult(null)}
        title="Great job!"
      />
      <IncorrectModal
        isOpen={submissionResult === "incorrect"}
        onClose={() => setSubmissionResult(null)}
      />
    </>
  );
}
