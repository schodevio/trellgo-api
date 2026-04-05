-- name: CreateCard :one
INSERT INTO cards (list_id, title, description, status, position)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetListCards :many
SELECT *
FROM cards
WHERE list_id = $1
ORDER BY position ASC;
