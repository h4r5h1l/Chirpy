-- +goose Up
CREATE TABLE refresh_tokens(
    token TEXT PRIMARY KEY,
    CREATED_AT TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP -- null if not revoked
);
-- +goose Down
DROP TABLE refresh_tokens;