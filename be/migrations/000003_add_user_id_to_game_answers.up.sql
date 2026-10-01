ALTER TABLE game_answers
ADD COLUMN IF NOT EXISTS user_id BIGINT REFERENCES users(id) ON DELETE CASCADE;

UPDATE game_answers AS ga
SET user_id = gs.user_id
FROM game_sessions AS gs
WHERE ga.game_session_id = gs.id
  AND ga.user_id IS NULL;

ALTER TABLE game_answers
ALTER COLUMN user_id SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_game_answers_user_vocabulary_unique
ON game_answers(user_id, vocabulary_id);