"use client";

import { useState } from "react";
import MainPage from "../components/MainPage";
import IncorrectModal from "../components/modal/IncorrectModal";
import SuccessModal from "../components/modal/SuccessModal";
import WordInput, { type GameWord } from "./WordInput";

type SubmissionResult = "correct" | "incorrect" | null;

export default function GamePage() {
  //TODO: use react-hook-form to handle form state and validation
  const [submissionResult, setSubmissionResult] =
    useState<SubmissionResult>(null);

  const word: GameWord = {
    id: 1,
    vocab: "brave",
    meaning:
      "giving you a benefit, a better position, or a higher chance of success",
    synonyms: "Beneficial, Useful, Helpful, Profitable",
    thai: "มีประโยชน์, เป็นประโยชน์, ได้เปรียบ",
    type_of_speech: "Adjective",
    status: "",
  };
  const degree = 1;
  const compliments = [
    "Great job!", 
    "You got it!", 
    "Awesome!", 
    "You're a genius!", 
    "You're Great!",
    "You're Amazing!",
    "You're Excellent!",
    "Perfect!",
  ];
  const title = compliments[Math.floor(Math.random() * compliments.length)];
  

  function handleSubmit(inputWord: string) {
    if (inputWord.toLowerCase() === word.vocab.toLowerCase()) {
      setSubmissionResult("correct");
      //TODO: update word status='success' to api
      return true;
    }
    setSubmissionResult("incorrect");
    return false;
  }

  function handleCloseModal() {
    setSubmissionResult(null);
  }
  function handleNextWord(){
    //TODO: get next word from API
  }

  return (
    <MainPage>
      <div>
        <WordInput word={word} degree={degree} onSubmit={handleSubmit} />
        <SuccessModal
          isOpen={submissionResult === "correct"}
          onSubmit={handleNextWord}
          onClose={handleCloseModal}
          title={title}
        />
        <IncorrectModal
          isOpen={submissionResult === "incorrect"}
          onClose={handleCloseModal}
        />
      </div>
    </MainPage>
  );
}
