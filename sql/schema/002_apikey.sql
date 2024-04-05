-- +goose Up
ALTER TABLE IF EXISTS users
ADD COLUMN API_KEY VARCHAR(64) NOT NULL DEFAULT encode(sha256(random()::text::bytea), 'hex'); 

-- +goose Down
ALTER TABLE users
DROP COLUMN API_KEY;

