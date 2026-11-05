"use client";

import { useEffect, useRef, useState } from "react";
import { useAuthStore } from "@/app/stores/auth.store";
import MainPage from "../components/ui/MainPage";
import Loading from "@/app/components/ui/Loading";
import LoadingModal from "@/app/components/modal/LoadingModal";
import { useModal } from "@/app/hooks/useModal";
import WordInput from "./WordInput";
import type { GameVocabResponse } from "@/app/type/GameType";
import { getVocab, submitVocabAnswer } from "@/app/services/game.service";

export default function GamePage() {
  const { showSuccessModal, showErrorModal } = useModal();
  const [gameResponse, setGameResponse] = useState<GameVocabResponse | null>(
    null,
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
      gameResponse.status === "success" ||
      inputWord.length !== gameResponse.length_of_word
    )
      return;

    submitting.current = true;
    setIsSubmitting(true);
    setError(null);
    try {
      const result = await submitVocabAnswer(gameResponse.id, inputWord);
      if (result.is_correct) {
        showSuccessModal({
          title: "Great job!",
          message: "That's the correct word!",
        });
      } else {
        showErrorModal({
          buttonText: "Try Again",
        });
      }
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
          <div
            className="grid min-h-screen place-items-center px-4 py-10"
            aria-busy="true"
          >
            <Loading
              title="Finding your next word…"
              message="A fresh challenge is on its way. Get ready to put your vocab to the test!"
            />
          </div>
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
      <LoadingModal
        isOpen={isSubmitting}
        title="Checking your answer…"
        message="Just a moment while we check your word."
      />
    </>
  );
}
