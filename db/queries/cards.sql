-- name: CreateCard :one
INSERT INTO cards (list_id, title, description, position)
VALUES ($1, $2, $3, $4)
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
SET title = $1, description = $2, updated_at = now()
WHERE id = $3
RETURNING *;

-- name: MoveCardByID :one
UPDATE cards
SET list_id = $1, position = $2, updated_at = now()
WHERE id = $3
RETURNING *;

-- name: ReorderCardsInList :exec
WITH ranked AS (
  SELECT
    id,
    ROW_NUMBER() OVER (
      PARTITION BY cards.list_id
      ORDER BY position ASC, updated_at DESC
    ) AS new_position
  FROM cards
  WHERE cards.list_id = $1
)
UPDATE cards
SET position = ranked.new_position
FROM ranked
WHERE cards.id = ranked.id;

-- name: DeleteCardByID :exec
DELETE FROM cards
WHERE id = $1;
