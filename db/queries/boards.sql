-- name: CreateBoard :one
INSERT INTO boards (name, user_id, created_at, updated_at)
VALUES ($1, $2, now(), now())
RETURNING *;

-- name: GetBoardsByUserID :many
SELECT *
FROM boards
WHERE user_id = $1
ORDER BY created_at DESC;
