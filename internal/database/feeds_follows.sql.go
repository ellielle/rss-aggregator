// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: feeds_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedsFollows = `-- name: CreateFeedsFollows :one
INSERT INTO feeds_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, user_id, feed_id
`

type CreateFeedsFollowsParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

func (q *Queries) CreateFeedsFollows(ctx context.Context, arg CreateFeedsFollowsParams) (FeedsFollow, error) {
	row := q.db.QueryRowContext(ctx, createFeedsFollows,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i FeedsFollow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
	)
	return i, err
}

const deleteFeedsFollows = `-- name: DeleteFeedsFollows :exec
DELETE FROM feeds_follows
WHERE id = $1
AND user_id = $2
`

type DeleteFeedsFollowsParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteFeedsFollows(ctx context.Context, arg DeleteFeedsFollowsParams) error {
	_, err := q.db.ExecContext(ctx, deleteFeedsFollows, arg.ID, arg.UserID)
	return err
}

const listFeedsFollows = `-- name: ListFeedsFollows :many
SELECT id, created_at, updated_at, user_id, feed_id FROM feeds_follows
WHERE user_id = $1
`

func (q *Queries) ListFeedsFollows(ctx context.Context, userID uuid.UUID) ([]FeedsFollow, error) {
	rows, err := q.db.QueryContext(ctx, listFeedsFollows, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedsFollow
	for rows.Next() {
		var i FeedsFollow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
