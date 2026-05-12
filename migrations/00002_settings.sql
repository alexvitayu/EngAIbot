-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_settings(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    language VARCHAR(100) NOT NULL,
    level VARCHAR(10) NOT NULL,
    topic VARCHAR(100) NOT NULL,
    interval_hours INT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_settings_user_id ON user_settings(user_id);
CREATE INDEX idx_user_settings_is_active ON user_settings(is_active);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_settings;
-- +goose StatementEnd
