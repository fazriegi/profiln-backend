-- name: ListPosts :many
SELECT *, COUNT(id) OVER () AS total_rows FROM posts 
ORDER BY updated_at DESC
OFFSET $1
LIMIT $2;

-- name: ListPostsByFollowing :many
SELECT p.*, COUNT(p.id) OVER () AS total_rows
FROM posts p
LEFT JOIN followings f ON p.user_id = f.follow_user_id
WHERE f.user_id = $1
ORDER BY p.updated_at DESC
OFFSET $2
LIMIT $3;

-- name: ListPopularPosts :many
SELECT *, COUNT(id) OVER () AS total_rows
FROM posts
ORDER BY 
	like_count DESC,
	repost_count DESC,
	comment_count DESC,
	updated_at DESC
OFFSET $1
LIMIT $2;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;