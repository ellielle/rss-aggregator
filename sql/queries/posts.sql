-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, feed_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetPostsByUser :many
SELECT posts.* FROM posts
JOIN feeds ON posts.feed_id = feeds.id
JOIN feeds_follows ON feeds.id = feeds_follows.feed_id
WHERE feeds_follows.user_id = $1
ORDER BY posts.updated_at DESC
LIMIT $2;
