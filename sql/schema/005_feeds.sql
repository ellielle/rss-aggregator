-- +goose Up
ALTER TABLE IF EXISTS feeds
ADD COLUMN last_fetched_at TIMESTAMP;

-- +goose Down
ALTER TABLE IF EXISTS feeds
DROP COLUMN IF EXISTS last_fetched_at;
