-- name: CreateUser :one
INSERT INTO users (id, provider, provider_id, username, email)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserExists :one
SELECT 1 FROM users WHERE id = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByProviderId :one
SELECT * FROM users WHERE provider = $1 AND provider_id = $2;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;
