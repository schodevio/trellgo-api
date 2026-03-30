-- name: CreateBoard :one
INSERT INTO boards (name, user_id, created_at, updated_at)
VALUES ($1, $2, now(), now())
RETURNING *;
