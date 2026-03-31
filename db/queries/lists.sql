-- name: CreateList :one
INSERT INTO lists (name, board_id, position)
VALUES ($1, $2, $3)
RETURNING *;
