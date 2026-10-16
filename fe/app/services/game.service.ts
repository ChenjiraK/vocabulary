import api from "./api";
import type { GameVocabResponse } from "@/app/type/GameType";

export async function getVocab(): Promise<GameVocabResponse | null> {
  const { data } = await api.get<{ vocabulary: GameVocabResponse | null }>(
    "/vocabulary/guess-word"
  );
  return data.vocabulary;
}

export async function submitVocabAnswer(id: number, word: string) {
  const { data } = await api.post<{ is_correct: boolean }>(
    "/vocabulary/guess-word/answer",
    { id, word }
  );
  return data;
}
