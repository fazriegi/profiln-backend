-- name: InsertUser :one
INSERT INTO users (
  email, password, full_name, verified_email
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2
WHERE id = $1
RETURNING *;