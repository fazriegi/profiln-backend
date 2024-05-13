-- name: InsertReportedPost :one
INSERT INTO reported_posts
(user_id, post_id, reason, message)
VALUES ($1, $2, $3, $4)
RETURNING *;