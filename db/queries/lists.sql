-- name: CreateList :one
INSERT INTO lists (name, board_id, position)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBoardLists :many
SELECT *
FROM lists
WHERE board_id = $1
ORDER BY position ASC;
