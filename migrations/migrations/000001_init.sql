-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS url
(
    alias        VARCHAR(10) PRIMARY KEY,
    original_url TEXT        NOT NULL UNIQUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT url_alias_length_check CHECK (char_length(alias) = 10),
    CONSTRAINT url_alias_format_check CHECK (alias ~ '^[A-Za-z0-9_]+$')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS url;
-- +goose StatementEnd