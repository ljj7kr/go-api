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


-- name: UpdateUser :execresult
UPDATE users
SET name = ?,
    email = ?
WHERE id = ?;


-- name: ListUsers :many
SELECT id,
       name,
       email,
       created_at
FROM users
ORDER BY id DESC LIMIT ?
OFFSET ?;


-- name: CountUsers :one
SELECT COUNT(*)
FROM users;


-- name: DeleteUser :execresult
DELETE
FROM users
WHERE id = ?;
