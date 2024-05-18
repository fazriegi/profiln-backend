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

-- name: UpdateVerifiedEmail :one
UPDATE users
SET verified_email = TRUE
FROM user_otps 
WHERE users.id = user_otps.user_id AND user_otps.otp = $1 AND users.email = $2
RETURNING users.id;

-- name: InsertOtp :one
INSERT INTO user_otps (
  user_id, otp
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetUserOtpByOtp :one
SELECT * 
FROM user_otps 
WHERE otp = $1 
LIMIT 1;

-- name: DeleteOtp :exec
DELETE FROM user_otps WHERE otp = $1;

-- name: GetUserOtpByEmail :one
SELECT *
FROM users
INNER JOIN user_otps ON users.id = user_otps.user_id
WHERE users.email = $1 AND users.verified_email = FALSE
LIMIT 1;