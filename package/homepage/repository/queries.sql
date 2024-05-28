-- name: ListNewestPosts :many
SELECT p.*, 
	u.id, u.full_name, u.avatar_url, u.bio, u.open_to_work, 
	COUNT(p.id) OVER () AS total_rows,
    CASE 
    	WHEN lp.user_id IS NOT NULL THEN TRUE 
    	ELSE FALSE 
  	END AS liked
FROM posts p
LEFT JOIN reported_posts rp ON p.id = rp.post_id AND rp.user_id = $1
LEFT JOIN users u ON p.user_id = u.id 
LEFT JOIN liked_posts lp ON p.id = lp.post_id
WHERE rp.post_id IS NULL
ORDER BY p.updated_at DESC
OFFSET $2
LIMIT $3;

-- name: ListPostsByFollowing :many
SELECT p.*,
	u.id, u.full_name, u.avatar_url, u.bio, u.open_to_work,
	COUNT(p.id) OVER () AS total_rows,
    CASE 
    	WHEN lp.user_id IS NOT NULL THEN TRUE 
    	ELSE FALSE 
  	END AS liked
FROM posts p
LEFT JOIN reported_posts rp ON p.id = rp.post_id AND rp.user_id = $1
LEFT JOIN users u ON p.user_id = u.id 
LEFT JOIN followings f ON p.user_id = f.follow_user_id
LEFT JOIN liked_posts lp ON p.id = lp.post_id
WHERE f.user_id = $1 AND rp.post_id IS NULL
ORDER BY p.updated_at DESC
OFFSET $2
LIMIT $3;

-- name: ListPopularPosts :many
SELECT p.*, 
	u.id, u.full_name, u.avatar_url, u.bio, u.open_to_work, 
	COUNT(p.id) OVER () AS total_rows,
    CASE
        WHEN p.created_at >= NOW() - INTERVAL '30 days' THEN true
        ELSE false
    END AS recent_post,
    CASE 
    	WHEN lp.user_id IS NOT NULL THEN TRUE 
    	ELSE FALSE 
  	END AS liked
FROM posts p
LEFT JOIN reported_posts rp ON p.id = rp.post_id AND rp.user_id = $1
LEFT JOIN users u ON p.user_id = u.id
LEFT JOIN liked_posts lp ON p.id = lp.post_id
WHERE rp.post_id IS NULL
ORDER BY
    recent_post DESC,
    (p.like_count + p.comment_count + p.repost_count) DESC
OFFSET $2
LIMIT $3;

-- name: GetFollowsRecommendationForUserId :many
SELECT u.id, u.full_name, u.avatar_url, u.bio, u.open_to_work, 
    COUNT(u.id) OVER () AS total_rows
FROM users u
JOIN (
    SELECT DISTINCT f.follow_user_id
    FROM followings f
    JOIN followings f2 ON f.user_id = f2.follow_user_id
    WHERE f2.user_id = $1
) AS users_reccomendation ON u.id = users_reccomendation.follow_user_id
WHERE u.id != $1
ORDER BY u.followers_count DESC
OFFSET $2
LIMIT $3;