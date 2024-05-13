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
SELECT *, COUNT(id) OVER () AS total_rows,
       CASE
           WHEN created_at >= NOW() - INTERVAL '30 days' THEN true
           ELSE false
       END AS recent_post
FROM posts
ORDER BY
    recent_post DESC,
    (like_count + comment_count + repost_count) DESC
OFFSET $1
LIMIT $2;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;