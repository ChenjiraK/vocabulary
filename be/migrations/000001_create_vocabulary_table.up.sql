CREATE TABLE IF NOT EXISTS vocabulary (
    id BIGSERIAL PRIMARY KEY,
    eng VARCHAR(255) NOT NULL,
    parts_of_speech VARCHAR(100) NOT NULL,
    thai VARCHAR(255) NOT NULL,
    meaning TEXT,
    synonyms TEXT,
    level VARCHAR(5) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
