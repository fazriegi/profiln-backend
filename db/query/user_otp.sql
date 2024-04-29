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