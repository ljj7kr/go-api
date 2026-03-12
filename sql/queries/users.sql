-- name: CreateUser :execresult
INSERT INTO users (name, email)
VALUES (?, ?);


-- name: GetUserByID :one
SELECT id,
       name,
       email,
       created_at
FROM users
WHERE id = ?;


-- name: GetUserByEmail :one
SELECT id,
       name,
       email,
       created_at
FROM users
WHERE email = ?;


-- name: ListUsers :many
SELECT id,
       name,
       email,
       created_at
FROM users
ORDER BY id DESC LIMIT ?
OFFSET ?;


-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = ?;