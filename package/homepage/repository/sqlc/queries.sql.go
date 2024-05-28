// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package homepage

import (
	"context"
	"database/sql"
)

const getFollowsRecommendationForUserId = `-- name: GetFollowsRecommendationForUserId :many
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
LIMIT $3
`

type GetFollowsRecommendationForUserIdParams struct {
	UserID sql.NullInt64
	Offset int32
	Limit  int32
}

type GetFollowsRecommendationForUserIdRow struct {
	ID         int64
	FullName   string
	AvatarUrl  sql.NullString
	Bio        sql.NullString
	OpenToWork sql.NullBool
	TotalRows  int64
}

func (q *Queries) GetFollowsRecommendationForUserId(ctx context.Context, arg GetFollowsRecommendationForUserIdParams) ([]GetFollowsRecommendationForUserIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getFollowsRecommendationForUserId, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFollowsRecommendationForUserIdRow
	for rows.Next() {
		var i GetFollowsRecommendationForUserIdRow
		if err := rows.Scan(
			&i.ID,
			&i.FullName,
			&i.AvatarUrl,
			&i.Bio,
			&i.OpenToWork,
			&i.TotalRows,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listNewestPosts = `-- name: ListNewestPosts :many
SELECT p.id, p.user_id, p.content, p.image_url, p.like_count, p.comment_count, p.repost_count, p.created_at, p.updated_at, p.title, p.visibility, 
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
LEFT JOIN reported_posts rp ON p.id = rp.post_id AND rp.user_id = $1
LEFT JOIN liked_posts lp ON p.id = lp.post_id AND lp.user_id = $1
LEFT JOIN reposted_posts rpp ON p.id = rpp.post_id AND rpp.user_id = $1
WHERE rp.post_id IS NULL AND p.visibility = 'public'
ORDER BY p.updated_at DESC
OFFSET $2
LIMIT $3
`

type ListNewestPostsParams struct {
	UserID sql.NullInt64
	Offset int32
	Limit  int32
}

type ListNewestPostsRow struct {
	ID           int64
	UserID       sql.NullInt64
	Content      sql.NullString
	ImageUrl     sql.NullString
	LikeCount    sql.NullInt32
	CommentCount sql.NullInt32
	RepostCount  sql.NullInt32
	CreatedAt    sql.NullTime
	UpdatedAt    sql.NullTime
	Title        string
	Visibility   string
	ID_2         sql.NullInt64
	FullName     sql.NullString
	AvatarUrl    sql.NullString
	Bio          sql.NullString
	OpenToWork   sql.NullBool
	TotalRows    int64
	Liked        bool
	Repost       bool
}

func (q *Queries) ListNewestPosts(ctx context.Context, arg ListNewestPostsParams) ([]ListNewestPostsRow, error) {
	rows, err := q.db.QueryContext(ctx, listNewestPosts, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListNewestPostsRow
	for rows.Next() {
		var i ListNewestPostsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Content,
			&i.ImageUrl,
			&i.LikeCount,
			&i.CommentCount,
			&i.RepostCount,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Visibility,
			&i.ID_2,
			&i.FullName,
			&i.AvatarUrl,
			&i.Bio,
			&i.OpenToWork,
			&i.TotalRows,
			&i.Liked,
			&i.Repost,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPopularPosts = `-- name: ListPopularPosts :many
SELECT p.id, p.user_id, p.content, p.image_url, p.like_count, p.comment_count, p.repost_count, p.created_at, p.updated_at, p.title, p.visibility, 
	u.id, u.full_name, u.avatar_url, u.bio, u.open_to_work, 
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
WHERE rp.post_id IS NULL AND p.visibility = 'public'
ORDER BY
    recent_post DESC,
    (p.like_count + p.comment_count + p.repost_count) DESC
OFFSET $2
LIMIT $3
`

type ListPopularPostsParams struct {
	UserID sql.NullInt64
	Offset int32
	Limit  int32
}

type ListPopularPostsRow struct {
	ID           int64
	UserID       sql.NullInt64
	Content      sql.NullString
	ImageUrl     sql.NullString
	LikeCount    sql.NullInt32
	CommentCount sql.NullInt32
	RepostCount  sql.NullInt32
	CreatedAt    sql.NullTime
	UpdatedAt    sql.NullTime
	Title        string
	Visibility   string
	ID_2         sql.NullInt64
	FullName     sql.NullString
	AvatarUrl    sql.NullString
	Bio          sql.NullString
	OpenToWork   sql.NullBool
	TotalRows    int64
	RecentPost   bool
	Liked        bool
	Repost       bool
}

func (q *Queries) ListPopularPosts(ctx context.Context, arg ListPopularPostsParams) ([]ListPopularPostsRow, error) {
	rows, err := q.db.QueryContext(ctx, listPopularPosts, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPopularPostsRow
	for rows.Next() {
		var i ListPopularPostsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Content,
			&i.ImageUrl,
			&i.LikeCount,
			&i.CommentCount,
			&i.RepostCount,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Visibility,
			&i.ID_2,
			&i.FullName,
			&i.AvatarUrl,
			&i.Bio,
			&i.OpenToWork,
			&i.TotalRows,
			&i.RecentPost,
			&i.Liked,
			&i.Repost,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPostsByFollowing = `-- name: ListPostsByFollowing :many
SELECT p.id, p.user_id, p.content, p.image_url, p.like_count, p.comment_count, p.repost_count, p.created_at, p.updated_at, p.title, p.visibility,
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
LEFT JOIN reported_posts rp ON p.id = rp.post_id AND rp.user_id = $1
LEFT JOIN followings f ON p.user_id = f.follow_user_id
LEFT JOIN liked_posts lp ON p.id = lp.post_id AND lp.user_id = $1
LEFT JOIN reposted_posts rpp ON p.id = rpp.post_id AND rpp.user_id = $1
WHERE f.user_id = $1 AND rp.post_id IS NULL
ORDER BY p.updated_at DESC
OFFSET $2
LIMIT $3
`

type ListPostsByFollowingParams struct {
	UserID sql.NullInt64
	Offset int32
	Limit  int32
}

type ListPostsByFollowingRow struct {
	ID           int64
	UserID       sql.NullInt64
	Content      sql.NullString
	ImageUrl     sql.NullString
	LikeCount    sql.NullInt32
	CommentCount sql.NullInt32
	RepostCount  sql.NullInt32
	CreatedAt    sql.NullTime
	UpdatedAt    sql.NullTime
	Title        string
	Visibility   string
	ID_2         sql.NullInt64
	FullName     sql.NullString
	AvatarUrl    sql.NullString
	Bio          sql.NullString
	OpenToWork   sql.NullBool
	TotalRows    int64
	Liked        bool
	Repost       bool
}

func (q *Queries) ListPostsByFollowing(ctx context.Context, arg ListPostsByFollowingParams) ([]ListPostsByFollowingRow, error) {
	rows, err := q.db.QueryContext(ctx, listPostsByFollowing, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPostsByFollowingRow
	for rows.Next() {
		var i ListPostsByFollowingRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Content,
			&i.ImageUrl,
			&i.LikeCount,
			&i.CommentCount,
			&i.RepostCount,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Visibility,
			&i.ID_2,
			&i.FullName,
			&i.AvatarUrl,
			&i.Bio,
			&i.OpenToWork,
			&i.TotalRows,
			&i.Liked,
			&i.Repost,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
