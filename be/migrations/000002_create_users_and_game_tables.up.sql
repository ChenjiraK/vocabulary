CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    display_name VARCHAR(100) NOT NULL,
    avatar_url TEXT,
    user_type VARCHAR(20) NOT NULL DEFAULT 'guest',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_user_type_check
        CHECK (user_type IN ('guest', 'registered'))
);

CREATE TABLE IF NOT EXISTS user_auth_identities (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(30) NOT NULL,
    provider_user_id VARCHAR(255),
    email VARCHAR(255),
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    password_hash TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT user_auth_identities_provider_check
        CHECK (provider IN ('password', 'google', 'facebook')),

    CONSTRAINT user_auth_identities_provider_identity_unique
        UNIQUE (provider, provider_user_id),

    CONSTRAINT user_auth_identities_email_provider_unique
        UNIQUE (provider, email)
);

CREATE TABLE IF NOT EXISTS game_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    total_questions INT NOT NULL DEFAULT 0,
    correct_answers INT NOT NULL DEFAULT 0,
    score INT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'in_progress',
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ,

    CONSTRAINT game_sessions_status_check
        CHECK (status IN ('in_progress', 'completed', 'abandoned'))
);

CREATE TABLE IF NOT EXISTS game_answers (
    id BIGSERIAL PRIMARY KEY,
    game_session_id BIGINT NOT NULL REFERENCES game_sessions(id) ON DELETE CASCADE,
    vocabulary_id BIGINT NOT NULL REFERENCES vocabulary(id) ON DELETE RESTRICT,
    question_text TEXT NOT NULL,
    correct_answer TEXT NOT NULL,
    user_answer TEXT,
    is_correct BOOLEAN NOT NULL DEFAULT FALSE,
    answered_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_auth_identities_user_id ON user_auth_identities(user_id);
CREATE INDEX IF NOT EXISTS idx_user_auth_identities_email ON user_auth_identities(email);
CREATE INDEX IF NOT EXISTS idx_game_sessions_user_id ON game_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_game_sessions_score ON game_sessions(score DESC);
CREATE INDEX IF NOT EXISTS idx_game_answers_session_id ON game_answers(game_session_id);
CREATE INDEX IF NOT EXISTS idx_game_answers_vocabulary_id ON game_answers(vocabulary_id);
