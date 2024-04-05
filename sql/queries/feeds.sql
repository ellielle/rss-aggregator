-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id) VALUES
($1, $2, $3, $4, $5),
("user_id", (SELECT id FROM users WHERE id = user_id))
RETURNING *;

