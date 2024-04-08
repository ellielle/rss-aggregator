-- +goose Up
ALTER TABLE IF EXISTS users
ADD COLUMN API_KEY VARCHAR(64) NOT NULL DEFAULT encode(sha256(random()::text::bytea), 'hex'); 
CREATE INDEX IF NOT EXISTS api_key_idx on users (api_key);

-- +goose Down
ALTER TABLE users
DROP COLUMN IF EXISTS API_KEY;
DROP INDEX IF EXISTS api_key_idx;
