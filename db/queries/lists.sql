-- name: CreateList :one
INSERT INTO lists (name, board_id, position)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBoardLists :many
SELECT *
FROM lists
WHERE board_id = $1
ORDER BY position ASC;

-- name: GetList :one
SELECT *
FROM lists
WHERE id = $1;

-- name: UpdateList :one
UPDATE lists
SET name = $1, position = $2, updated_at = now()
WHERE id = $3
RETURNING *;
