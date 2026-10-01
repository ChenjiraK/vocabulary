DROP INDEX IF EXISTS idx_game_answers_user_vocabulary_unique;

ALTER TABLE game_answers
DROP COLUMN IF EXISTS user_id;