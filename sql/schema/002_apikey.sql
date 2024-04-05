-- +goose Up
ALTER TABLE users
ADD COLUMN API_KEY TEXT NOT NULL UNIQUE (
				SELECT encode(sha256(random()::text::bytea), 'hex')
);

-- +goose Down
ALTER TABLE users
DROP COLUMN API_KEY;
