-- name: CreateUser :exec
INSERT INTO users (email, password_hash, created_at, updated_at)
VALUES ($1, $2, now(), now())
RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;
