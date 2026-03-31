-- name: CreateBoard :one
INSERT INTO boards (name, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserBoards :many
SELECT *
FROM boards
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetUserBoardByID :one
SELECT *
FROM boards
WHERE id = $1 AND user_id = $2;

-- name: UpdateUserBoardByID :one
UPDATE boards
SET name = $1, updated_at = now()
WHERE id = $2 AND user_id = $3
RETURNING *;

-- name: DeleteUserBoardByID :exec
DELETE FROM boards
WHERE id = $1 AND user_id = $2;
