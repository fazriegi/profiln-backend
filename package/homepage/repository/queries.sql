-- name: ListPosts :many
SELECT p.*, COUNT(p.id) OVER () AS total_rows 
FROM posts p
LEFT JOIN reported_posts rp ON p.id = rp.post_id AND rp.user_id = $1
WHERE rp.post_id IS NULL
ORDER BY updated_at DESC
OFFSET $2
LIMIT $3;

-- name: ListPostsByFollowing :many
SELECT p.*, COUNT(p.id) OVER () AS total_rows
FROM posts p
LEFT JOIN reported_posts rp ON p.id = rp.post_id AND rp.user_id = $1
LEFT JOIN followings f ON p.user_id = f.follow_user_id
WHERE f.user_id = $1 AND rp.post_id IS NULL
ORDER BY p.updated_at DESC
OFFSET $2
LIMIT $3;

-- name: ListPopularPosts :many
SELECT p.*, COUNT(p.id) OVER () AS total_rows,
       CASE
           WHEN p.created_at >= NOW() - INTERVAL '30 days' THEN true
           ELSE false
       END AS recent_post
FROM posts p
LEFT JOIN reported_posts rp ON p.id = rp.post_id AND rp.user_id = $1
WHERE rp.post_id IS NULL
ORDER BY
    recent_post DESC,
    (p.like_count + p.comment_count + p.repost_count) DESC
OFFSET $2
LIMIT $3;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;