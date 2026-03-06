-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS url
(
    alias        VARCHAR(10) PRIMARY KEY,
    original_url TEXT        NOT NULL UNIQUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS url;
-- +goose StatementEnd
