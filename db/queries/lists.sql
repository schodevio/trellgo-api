-- name: CreateList :one
INSERT INTO lists (name, board_id, position)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBoardLists :many
SELECT *
FROM lists
WHERE board_id = $1
ORDER BY position ASC;

-- name: GetUserListByID :one
SELECT lists.*
FROM lists
INNER JOIN boards ON lists.board_id = boards.id
WHERE lists.id = $1 AND boards.user_id = $2;

-- name: GetListByID :one
SELECT *
FROM lists
WHERE id = $1;

-- name: UpdateListByID :one
UPDATE lists
SET name = $1, updated_at = now()
WHERE id = $2
RETURNING *;

-- name: DeleteListByID :exec
DELETE FROM lists
WHERE id = $1;
