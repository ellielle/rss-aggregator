-- +goose Up
CREATE TABLE IF NOT EXISTS feeds_follows (
				id UUID PRIMARY KEY,
				created_at TIMESTAMP NOT NULL,
				updated_at TIMESTAMP NOT NULL,
				user_id UUID NOT NULL,
				feed_id UUID NOT NULL,
				CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE, 
				CONSTRAINT fk_feeds FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS feeds_follows;

