-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS phrase(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    target_language VARCHAR(50) NOT NULL,
    level VARCHAR(10) NOT NULL,
    topic VARCHAR(100) NOT NULL,
    in_language_text TEXT NOT NULL,
    in_russian_text TEXT NOT NULL,
    generated_by VARCHAR(50) DEFAULT 'groq',
    usage_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_phrases_language_level_topic ON phrase (target_language, level, topic);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS phrase;
-- +goose StatementEnd
