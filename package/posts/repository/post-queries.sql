-- name: InsertReportedPost :one
INSERT INTO reported_posts
(user_id, post_id, reason, message)
SELECT @user_id::bigint, @post_id::bigint, UNNEST(@reason::varchar(15)[]), @message::text
RETURNING *;

-- name: GetDetailPost :one
SELECT p.*, 
    pu.id, pu.avatar_url, pu.full_name, pu.bio, pu.open_to_work,
	ARRAY_AGG(pi.url ORDER BY pi.index ASC) FILTER (WHERE pi.url IS NOT NULL) AS image_urls,
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
LEFT JOIN post_images pi ON p.id = pi.post_id
WHERE p.id = $1
GROUP BY 
    p.id, pu.id, lp.user_id, rpp.user_id;

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
SET like_count = GREATEST(like_count + @value::smallint, 0)
WHERE id = @id::bigint
RETURNING id, like_count;

-- name: ListNewestPostsByUserId :many
SELECT p.*, 
	u.id, u.full_name, u.avatar_url, u.bio, u.open_to_work,
	ARRAY_AGG(pi.url) FILTER (WHERE pi.url IS NOT NULL) AS image_urls,
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
LEFT JOIN post_images pi ON p.id = pi.post_id
WHERE p.user_id = $1
GROUP BY 
    p.id, u.id, lp.user_id, rpp.user_id
ORDER BY p.updated_at DESC
OFFSET $2
LIMIT $3;

-- name: ListLikedPostsByUserId :many
SELECT p.*, 
	u.id, u.full_name, u.avatar_url, u.bio, u.open_to_work,
	ARRAY_AGG(pi.url) FILTER (WHERE pi.url IS NOT NULL) AS image_urls,
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
LEFT JOIN post_images pi ON p.id = pi.post_id
WHERE lp.user_id = $1
GROUP BY 
    p.id, u.id, lp.user_id, rpp.user_id
ORDER BY p.updated_at DESC
OFFSET $2
LIMIT $3;

-- name: ListRepostedPostsByUserId :many
SELECT p.*, 
	u.id, u.full_name, u.avatar_url, u.bio, u.open_to_work,
	ARRAY_AGG(pi.url) FILTER (WHERE pi.url IS NOT NULL) AS image_urls,
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
LEFT JOIN post_images pi ON p.id = pi.post_id
WHERE rpp.user_id = $1
GROUP BY 
    p.id, u.id, lp.user_id, rpp.user_id
ORDER BY p.updated_at DESC
OFFSET $2
LIMIT $3;

-- name: InsertPost :one
INSERT INTO posts
(user_id, title, content, visibility)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPostById :one
SELECT * FROM posts
WHERE id = $1
LIMIT 1;

-- name: UpdatePost :exec
UPDATE posts
SET title = @title::text,
    content = @content::text,
    visibility = @visibility::varchar(10)
WHERE id = @id::bigint AND user_id = @user_id::bigint;

-- name: DeletePostById :exec
DELETE FROM posts
WHERE id = @id::bigint;

-- name: BatchInsertPostImages :many
INSERT INTO post_images
	(post_id, url, index)
SELECT @post_id::bigint, UNNEST(@url::TEXT[]), UNNEST(@index::smallint[])
RETURNING *;

-- name: GetPostImagesUrl :many
SELECT url FROM post_images
WHERE post_id = @post_id::bigint;

-- name: BatchDeletePostImagesByPost :exec
DELETE FROM post_images
WHERE post_id = @post_id::bigint;

-- name: BatchDeleteReportedPostsByPost :exec
DELETE FROM reported_posts
WHERE post_id = @post_id::bigint;

-- name: BatchDeleteLikedPostByPost :exec
DELETE FROM liked_posts
WHERE post_id = @post_id::bigint;

-- name: BatchDeleteRepostedPostByPost :exec
DELETE FROM reposted_posts
WHERE post_id = @post_id::bigint;

-- name: BatchDeletePostCommentsByPost :exec
DELETE FROM post_comments
WHERE post_id = @post_id::bigint;

-- name: BatchDeletePostCommentRepliesByPost :exec
WITH post_comments AS (
	SELECT id 
	FROM post_comments 
	WHERE post_id = @post_id::bigint
)
DELETE FROM post_comment_replies
WHERE post_comment_id IN (SELECT id FROM post_comments);

-- name: InsertLikedPost :one
INSERT INTO liked_posts (user_id, post_id)
VALUES (@user_id::bigint, @post_id::bigint)
ON CONFLICT (user_id, post_id) DO NOTHING
RETURNING id;

-- name: DeleteLikedPost :one
DELETE FROM liked_posts
WHERE user_id = @user_id::bigint AND post_id = @post_id::bigint
RETURNING id;

-- name: UpdatePostRepostCount :one
UPDATE posts
SET repost_count = GREATEST(repost_count + @value::smallint, 0)
WHERE id = @id::bigint
RETURNING id, repost_count;

-- name: InsertRepostedPost :one
INSERT INTO reposted_posts (user_id, post_id)
VALUES (@user_id::bigint, @post_id::bigint)
ON CONFLICT (user_id, post_id) DO NOTHING
RETURNING id;

-- name: DeleteRepostedPost :one
DELETE FROM reposted_posts
WHERE user_id = @user_id::bigint AND post_id = @post_id::bigint
RETURNING id;

-- name: CountPostImages :one
SELECT COUNT(*) AS count
FROM post_images
WHERE post_id = @post_id::bigint;