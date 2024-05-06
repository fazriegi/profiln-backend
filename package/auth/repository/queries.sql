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

-- name: UpdateVerifiedEmailByOTP :exec
UPDATE users
SET verified_email = TRUE
FROM user_otps 
WHERE users.id = user_otps.id AND user_otps.otp = $1
RETURNING *;

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