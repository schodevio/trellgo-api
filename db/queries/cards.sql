-- name: CreateCard :one
INSERT INTO cards (list_id, title, description, status, position)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetListCards :many
SELECT *
FROM cards
WHERE list_id = $1
ORDER BY position ASC;

-- name: GetCardByID :one
SELECT *
FROM cards
WHERE id = $1;

-- name: UpdateCardByID :one
UPDATE cards
SET title = $1, description = $2, status = $3, updated_at = now()
WHERE id = $4
RETURNING *;

-- name: DeleteCardByID :exec
DELETE FROM cards
WHERE id = $1;
