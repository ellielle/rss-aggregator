-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, feed_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetPostsByUser :many
SELECT * FROM posts
JOIN users ON users.id = $1
ORDER BY posts.updated_at DESC
LIMIT $2;
