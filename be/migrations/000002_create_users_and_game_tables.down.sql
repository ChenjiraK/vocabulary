DROP INDEX IF EXISTS idx_game_answers_vocabulary_id;
DROP INDEX IF EXISTS idx_game_answers_session_id;
DROP INDEX IF EXISTS idx_game_sessions_score;
DROP INDEX IF EXISTS idx_game_sessions_user_id;
DROP INDEX IF EXISTS idx_user_auth_identities_email;
DROP INDEX IF EXISTS idx_user_auth_identities_user_id;

DROP TABLE IF EXISTS game_answers;
DROP TABLE IF EXISTS game_sessions;
DROP TABLE IF EXISTS user_auth_identities;
DROP TABLE IF EXISTS users;
