-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    tg_user_id BIGINT NOT NULL,
    user_name varchar(100),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    chat_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_tg_user_id ON users (tg_user_id);
CREATE INDEX IF NOT EXISTS idx_users_chat_id ON users (chat_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
