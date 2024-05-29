-- name: InsertReportedPost :one
INSERT INTO reported_posts
(user_id, post_id, reason, message)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetDetailPost :one
SELECT p.*, 
    pu.id, pu.avatar_url, pu.full_name, pu.bio, pu.open_to_work,
    CASE 
    	WHEN lp.user_id IS NOT NULL THEN TRUE 
    	ELSE FALSE 
  	END AS liked,
	CASE 
    	WHEN rpp.user_id IS NOT NULL THEN TRUE 
    	ELSE FALSE 
  	END AS repost
FROM posts p
JOIN users pu ON p.user_id = pu.id
LEFT JOIN liked_posts lp ON p.id = lp.post_id AND lp.user_id = $2
LEFT JOIN reposted_posts rpp ON p.id = rpp.post_id AND rpp.user_id = $2
WHERE p.id = $1;

-- name: GetPostComments :many
SELECT pc.*,
    pcu.id, pcu.avatar_url, pcu.full_name, pcu.bio, pcu.open_to_work,
    COUNT(pc.id) OVER () AS total_rows
FROM post_comments pc 
LEFT JOIN users pcu ON pc.user_id = pcu.id
WHERE pc.post_id = $1
ORDER BY pc.updated_at DESC
OFFSET $2
LIMIT $3;

-- name: GetPostCommentReplies :many
SELECT pcr.*, 
    pcr_user.id, pcr_user.avatar_url, pcr_user.full_name, pcr_user.bio, pcr_user.open_to_work,
    COUNT(pcr.id) OVER () AS total_rows
FROM post_comment_replies pcr 
LEFT JOIN users pcr_user ON pcr.user_id = pcr_user.id
LEFT JOIN post_comments pc ON pc.id = pcr.post_comment_id
WHERE pc.post_id = $1 AND pcr.post_comment_id = $2
ORDER BY pcr.updated_at DESC
OFFSET $3
LIMIT $4;

-- name: LockPostForUpdate :one
SELECT 1
FROM posts
WHERE id = $1
FOR UPDATE;

-- name: UpdatePostLikeCount :one
UPDATE posts
SET like_count = like_count + 1
WHERE id = $1
RETURNING id, like_count;

-- name: ListNewestPostsByUserId :many
SELECT p.*, 
	u.id, u.full_name, u.avatar_url, u.bio, u.open_to_work, 
	COUNT(p.id) OVER () AS total_rows,
    CASE 
    	WHEN lp.user_id IS NOT NULL THEN TRUE 
    	ELSE FALSE 
  	END AS liked,
	CASE 
    	WHEN rpp.user_id IS NOT NULL THEN TRUE 
    	ELSE FALSE 
  	END AS repost
FROM posts p
LEFT JOIN users u ON p.user_id = u.id
LEFT JOIN liked_posts lp ON p.id = lp.post_id AND lp.user_id = $1
LEFT JOIN reposted_posts rpp ON p.id = rpp.post_id AND rpp.user_id = $1
WHERE p.user_id = $1
ORDER BY p.updated_at DESC
OFFSET $2
LIMIT $3;

-- name: ListLikedPostsByUserId :many
SELECT p.*, 
	u.id, u.full_name, u.avatar_url, u.bio, u.open_to_work, 
	COUNT(p.id) OVER () AS total_rows,
    CASE 
    	WHEN lp.user_id IS NOT NULL THEN TRUE 
    	ELSE FALSE 
  	END AS liked,
	CASE 
    	WHEN rpp.user_id IS NOT NULL THEN TRUE 
    	ELSE FALSE 
  	END AS repost
FROM posts p
LEFT JOIN users u ON p.user_id = u.id
LEFT JOIN liked_posts lp ON p.id = lp.post_id AND lp.user_id = $1
LEFT JOIN reposted_posts rpp ON p.id = rpp.post_id AND rpp.user_id = $1
WHERE lp.user_id = $1
ORDER BY p.updated_at DESC
OFFSET $2
LIMIT $3;

-- name: ListRepostedPostsByUserId :many
SELECT p.*, 
	u.id, u.full_name, u.avatar_url, u.bio, u.open_to_work, 
	COUNT(p.id) OVER () AS total_rows,
    CASE 
    	WHEN lp.user_id IS NOT NULL THEN TRUE 
    	ELSE FALSE 
  	END AS liked,
	CASE 
    	WHEN rpp.user_id IS NOT NULL THEN TRUE 
    	ELSE FALSE 
  	END AS repost
FROM posts p
LEFT JOIN users u ON p.user_id = u.id
LEFT JOIN liked_posts lp ON p.id = lp.post_id AND lp.user_id = $1
LEFT JOIN reposted_posts rpp ON p.id = rpp.post_id AND rpp.user_id = $1
WHERE rpp.user_id = $1
ORDER BY p.updated_at DESC
OFFSET $2
LIMIT $3;

-- name: InsertPost :one
INSERT INTO posts
(user_id, title, content, visibility)
VALUES ($1, $2, $3, $4)
RETURNING *;