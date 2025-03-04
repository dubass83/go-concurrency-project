-- name: InsertUser :one
INSERT INTO users (
  email, first_name, last_name, password, user_active
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetOneUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetAllUsers :many
SELECT * FROM users
ORDER by last_name
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET
  first_name = COALESCE(sqlc.narg('first_name'), first_name),
  last_name = COALESCE(sqlc.narg('last_name'), last_name),
  email = COALESCE(sqlc.narg('email'), email),
  user_active = COALESCE(sqlc.narg('user_active'), user_active),
  password = COALESCE(sqlc.narg('password'), password),
  updated_at = COALESCE(sqlc.narg('updated_at'), updated_at)
WHERE
  id = sqlc.arg('id')
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE id = $1;
