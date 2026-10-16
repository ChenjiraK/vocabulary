export type GameVocabResponse = {
  id: number;
  status?: "" | "success" | "failed";
  length_of_word: number;
  parts_of_speech: string;
  thai: string;
  start_letter: string | "A";
  end_letter: string | "Z";
  meaning: string;
  synonyms: string;
};
