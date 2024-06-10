-- name: ListNewestPosts :many
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
LEFT JOIN reported_posts rp ON p.id = rp.post_id AND rp.user_id = $1
LEFT JOIN liked_posts lp ON p.id = lp.post_id AND lp.user_id = $1
LEFT JOIN reposted_posts rpp ON p.id = rpp.post_id AND rpp.user_id = $1
LEFT JOIN post_images pi ON p.id = pi.post_id
WHERE rp.post_id IS NULL AND p.visibility = 'public'
GROUP BY 
    p.id, u.id, lp.user_id, rpp.user_id
ORDER BY p.created_at DESC
OFFSET $2
LIMIT $3;

-- name: ListPostsByFollowing :many
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
LEFT JOIN reported_posts rp ON p.id = rp.post_id AND rp.user_id = $1
LEFT JOIN followings f ON p.user_id = f.follow_user_id
LEFT JOIN liked_posts lp ON p.id = lp.post_id AND lp.user_id = $1
LEFT JOIN reposted_posts rpp ON p.id = rpp.post_id AND rpp.user_id = $1
LEFT JOIN post_images pi ON p.id = pi.post_id
WHERE f.user_id = $1 AND rp.post_id IS NULL
GROUP BY 
    p.id, u.id, lp.user_id, rpp.user_id
ORDER BY p.created_at DESC
OFFSET $2
LIMIT $3;

-- name: ListPopularPosts :many
SELECT p.*, 
	u.id, u.full_name, u.avatar_url, u.bio, u.open_to_work,
	ARRAY_AGG(pi.url) FILTER (WHERE pi.url IS NOT NULL) AS image_urls,
	COUNT(p.id) OVER () AS total_rows,
    CASE
        WHEN p.created_at >= NOW() - INTERVAL '30 days' THEN true
        ELSE false
    END AS recent_post,
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
LEFT JOIN reported_posts rp ON p.id = rp.post_id AND rp.user_id = $1
LEFT JOIN liked_posts lp ON p.id = lp.post_id AND lp.user_id = $1
LEFT JOIN reposted_posts rpp ON p.id = rpp.post_id AND rpp.user_id = $1
LEFT JOIN post_images pi ON p.id = pi.post_id
WHERE rp.post_id IS NULL AND p.visibility = 'public'
GROUP BY 
    p.id, u.id, lp.user_id, rpp.user_id
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