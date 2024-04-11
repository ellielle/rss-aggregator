-- +goose Up
CREATE TABLE IF NOT EXISTS posts (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	title TEXT NOT NULL,
	url TEXT NOT NULL UNIQUE,
	description TEXT,
	feed_id UUID NOT NULL,
	CONSTRAINT fk_feeds
	FOREIGN KEY (feed_id)
	REFERENCES feeds(id)
	ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS posts;
