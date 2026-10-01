export type VocabularyLevel = "A1" | "A2" | "B1" | "B2";

export type VocabularyItem = {
  id: number;
  word: string;
  meaning: string;
  partOfSpeech: string;
  level: VocabularyLevel;
  example: string;
  lastReviewed: string;
};

export const vocabulary: VocabularyItem[] = [
  {
    id: 1,
    word: "curious",
    meaning: "eager to learn or know",
    partOfSpeech: "adjective",
    level: "A2",
    example: "She is curious about how languages work.",
    lastReviewed: "2026-08-10",
  },
  {
    id: 2,
    word: "steady",
    meaning: "firm, controlled, and reliable",
    partOfSpeech: "adjective",
    level: "B1",
    example: "A steady routine makes practice easier.",
    lastReviewed: "2026-08-11",
  },
  {
    id: 3,
    word: "bright",
    meaning: "full of light or intelligence",
    partOfSpeech: "adjective",
    level: "A1",
    example: "The room felt bright in the morning.",
    lastReviewed: "2026-08-09",
  },
  {
    id: 4,
    word: "approach",
    meaning: "a way of dealing with something",
    partOfSpeech: "noun",
    level: "B1",
    example: "Her approach to study is calm and consistent.",
    lastReviewed: "2026-08-12",
  },
  {
    id: 5,
    word: "expand",
    meaning: "to become larger or make something larger",
    partOfSpeech: "verb",
    level: "A2",
    example: "Reading every day can expand your vocabulary.",
    lastReviewed: "2026-08-08",
  },
  {
    id: 6,
    word: "precise",
    meaning: "exact and accurate",
    partOfSpeech: "adjective",
    level: "B2",
    example: "A precise definition helps you remember the word.",
    lastReviewed: "2026-08-13",
  },
];

export async function fetchVocabulary() {
  await new Promise((resolve) => setTimeout(resolve, 250));
  return vocabulary;
}
