// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: post-queries.sql

package db

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

const batchDeleteLikedPostByPost = `-- name: BatchDeleteLikedPostByPost :exec
DELETE FROM liked_posts
WHERE post_id = $1::bigint
`

func (q *Queries) BatchDeleteLikedPostByPost(ctx context.Context, postID int64) error {
	_, err := q.db.ExecContext(ctx, batchDeleteLikedPostByPost, postID)
	return err
}

const batchDeletePostCommentRepliesByPost = `-- name: BatchDeletePostCommentRepliesByPost :exec
WITH post_comments AS (
	SELECT id 
	FROM post_comments 
	WHERE post_id = $1::bigint
)
DELETE FROM post_comment_replies
WHERE post_comment_id IN (SELECT id FROM post_comments)
`

func (q *Queries) BatchDeletePostCommentRepliesByPost(ctx context.Context, postID int64) error {
	_, err := q.db.ExecContext(ctx, batchDeletePostCommentRepliesByPost, postID)
	return err
}

const batchDeletePostCommentsByPost = `-- name: BatchDeletePostCommentsByPost :exec
DELETE FROM post_comments
WHERE post_id = $1::bigint
`

func (q *Queries) BatchDeletePostCommentsByPost(ctx context.Context, postID int64) error {
	_, err := q.db.ExecContext(ctx, batchDeletePostCommentsByPost, postID)
	return err
}

const batchDeletePostImagesByPost = `-- name: BatchDeletePostImagesByPost :exec
DELETE FROM post_images
WHERE post_id = $1::bigint
`

func (q *Queries) BatchDeletePostImagesByPost(ctx context.Context, postID int64) error {
	_, err := q.db.ExecContext(ctx, batchDeletePostImagesByPost, postID)
	return err
}

const batchDeleteReportedPostsByPost = `-- name: BatchDeleteReportedPostsByPost :exec
DELETE FROM reported_posts
WHERE post_id = $1::bigint
`

func (q *Queries) BatchDeleteReportedPostsByPost(ctx context.Context, postID int64) error {
	_, err := q.db.ExecContext(ctx, batchDeleteReportedPostsByPost, postID)
	return err
}

const batchDeleteRepostedPostByPost = `-- name: BatchDeleteRepostedPostByPost :exec
DELETE FROM reposted_posts
WHERE post_id = $1::bigint
`

func (q *Queries) BatchDeleteRepostedPostByPost(ctx context.Context, postID int64) error {
	_, err := q.db.ExecContext(ctx, batchDeleteRepostedPostByPost, postID)
	return err
}

const batchInsertPostImages = `-- name: BatchInsertPostImages :many
INSERT INTO post_images
	(post_id, url, index)
SELECT $1::bigint, UNNEST($2::TEXT[]), UNNEST($3::smallint[])
RETURNING id, post_id, url, index
`

type BatchInsertPostImagesParams struct {
	PostID int64
	Url    []string
	Index  []int16
}

func (q *Queries) BatchInsertPostImages(ctx context.Context, arg BatchInsertPostImagesParams) ([]PostImage, error) {
	rows, err := q.db.QueryContext(ctx, batchInsertPostImages, arg.PostID, pq.Array(arg.Url), pq.Array(arg.Index))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PostImage
	for rows.Next() {
		var i PostImage
		if err := rows.Scan(
			&i.ID,
			&i.PostID,
			&i.Url,
			&i.Index,
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

const countPostImages = `-- name: CountPostImages :one
SELECT COUNT(*) AS count
FROM post_images
WHERE post_id = $1::bigint
`

func (q *Queries) CountPostImages(ctx context.Context, postID int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, countPostImages, postID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const deleteLikedPost = `-- name: DeleteLikedPost :one
DELETE FROM liked_posts
WHERE user_id = $1::bigint AND post_id = $2::bigint
RETURNING id
`

type DeleteLikedPostParams struct {
	UserID int64
	PostID int64
}

func (q *Queries) DeleteLikedPost(ctx context.Context, arg DeleteLikedPostParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, deleteLikedPost, arg.UserID, arg.PostID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const deletePostById = `-- name: DeletePostById :exec
DELETE FROM posts
WHERE id = $1::bigint
`

func (q *Queries) DeletePostById(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deletePostById, id)
	return err
}

const deleteRepostedPost = `-- name: DeleteRepostedPost :one
DELETE FROM reposted_posts
WHERE user_id = $1::bigint AND post_id = $2::bigint
RETURNING id
`

type DeleteRepostedPostParams struct {
	UserID int64
	PostID int64
}

func (q *Queries) DeleteRepostedPost(ctx context.Context, arg DeleteRepostedPostParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, deleteRepostedPost, arg.UserID, arg.PostID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getDetailPost = `-- name: GetDetailPost :one
SELECT p.id, p.user_id, p.content, p.like_count, p.comment_count, p.repost_count, p.created_at, p.updated_at, p.title, p.visibility, 
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
    p.id, pu.id, lp.user_id, rpp.user_id
`

type GetDetailPostParams struct {
	ID     int64
	UserID sql.NullInt64
}

type GetDetailPostRow struct {
	ID           int64
	UserID       sql.NullInt64
	Content      sql.NullString
	LikeCount    sql.NullInt32
	CommentCount sql.NullInt32
	RepostCount  sql.NullInt32
	CreatedAt    sql.NullTime
	UpdatedAt    sql.NullTime
	Title        string
	Visibility   string
	ID_2         int64
	AvatarUrl    sql.NullString
	FullName     string
	Bio          sql.NullString
	OpenToWork   sql.NullBool
	ImageUrls    interface{}
	Liked        bool
	Repost       bool
}

func (q *Queries) GetDetailPost(ctx context.Context, arg GetDetailPostParams) (GetDetailPostRow, error) {
	row := q.db.QueryRowContext(ctx, getDetailPost, arg.ID, arg.UserID)
	var i GetDetailPostRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Content,
		&i.LikeCount,
		&i.CommentCount,
		&i.RepostCount,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Visibility,
		&i.ID_2,
		&i.AvatarUrl,
		&i.FullName,
		&i.Bio,
		&i.OpenToWork,
		&i.ImageUrls,
		&i.Liked,
		&i.Repost,
	)
	return i, err
}

const getPostById = `-- name: GetPostById :one
SELECT id, user_id, content, like_count, comment_count, repost_count, created_at, updated_at, title, visibility FROM posts
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetPostById(ctx context.Context, id int64) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostById, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Content,
		&i.LikeCount,
		&i.CommentCount,
		&i.RepostCount,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Visibility,
	)
	return i, err
}

const getPostCommentReplies = `-- name: GetPostCommentReplies :many
SELECT pcr.id, pcr.user_id, pcr.post_comment_id, pcr.content, pcr.image_url, pcr.like_count, pcr.is_post_author, pcr.created_at, pcr.updated_at, 
    pcr_user.id, pcr_user.avatar_url, pcr_user.full_name, pcr_user.bio, pcr_user.open_to_work,
    COUNT(pcr.id) OVER () AS total_rows
FROM post_comment_replies pcr 
LEFT JOIN users pcr_user ON pcr.user_id = pcr_user.id
LEFT JOIN post_comments pc ON pc.id = pcr.post_comment_id
WHERE pc.post_id = $1 AND pcr.post_comment_id = $2
ORDER BY pcr.updated_at DESC
OFFSET $3
LIMIT $4
`

type GetPostCommentRepliesParams struct {
	PostID        sql.NullInt64
	PostCommentID sql.NullInt64
	Offset        int32
	Limit         int32
}

type GetPostCommentRepliesRow struct {
	ID            int64
	UserID        sql.NullInt64
	PostCommentID sql.NullInt64
	Content       sql.NullString
	ImageUrl      sql.NullString
	LikeCount     sql.NullInt32
	IsPostAuthor  sql.NullBool
	CreatedAt     sql.NullTime
	UpdatedAt     sql.NullTime
	ID_2          sql.NullInt64
	AvatarUrl     sql.NullString
	FullName      sql.NullString
	Bio           sql.NullString
	OpenToWork    sql.NullBool
	TotalRows     int64
}

func (q *Queries) GetPostCommentReplies(ctx context.Context, arg GetPostCommentRepliesParams) ([]GetPostCommentRepliesRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostCommentReplies,
		arg.PostID,
		arg.PostCommentID,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostCommentRepliesRow
	for rows.Next() {
		var i GetPostCommentRepliesRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.PostCommentID,
			&i.Content,
			&i.ImageUrl,
			&i.LikeCount,
			&i.IsPostAuthor,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID_2,
			&i.AvatarUrl,
			&i.FullName,
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

const getPostComments = `-- name: GetPostComments :many
SELECT pc.id, pc.user_id, pc.post_id, pc.content, pc.image_url, pc.like_count, pc.reply_count, pc.is_post_author, pc.created_at, pc.updated_at,
    pcu.id, pcu.avatar_url, pcu.full_name, pcu.bio, pcu.open_to_work,
    COUNT(pc.id) OVER () AS total_rows
FROM post_comments pc 
LEFT JOIN users pcu ON pc.user_id = pcu.id
WHERE pc.post_id = $1
ORDER BY pc.updated_at DESC
OFFSET $2
LIMIT $3
`

type GetPostCommentsParams struct {
	PostID sql.NullInt64
	Offset int32
	Limit  int32
}

type GetPostCommentsRow struct {
	ID           int64
	UserID       sql.NullInt64
	PostID       sql.NullInt64
	Content      sql.NullString
	ImageUrl     sql.NullString
	LikeCount    sql.NullInt32
	ReplyCount   sql.NullInt32
	IsPostAuthor sql.NullBool
	CreatedAt    sql.NullTime
	UpdatedAt    sql.NullTime
	ID_2         sql.NullInt64
	AvatarUrl    sql.NullString
	FullName     sql.NullString
	Bio          sql.NullString
	OpenToWork   sql.NullBool
	TotalRows    int64
}

func (q *Queries) GetPostComments(ctx context.Context, arg GetPostCommentsParams) ([]GetPostCommentsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostComments, arg.PostID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostCommentsRow
	for rows.Next() {
		var i GetPostCommentsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.PostID,
			&i.Content,
			&i.ImageUrl,
			&i.LikeCount,
			&i.ReplyCount,
			&i.IsPostAuthor,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID_2,
			&i.AvatarUrl,
			&i.FullName,
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

const getPostImagesUrl = `-- name: GetPostImagesUrl :many
SELECT url FROM post_images
WHERE post_id = $1::bigint
`

func (q *Queries) GetPostImagesUrl(ctx context.Context, postID int64) ([]sql.NullString, error) {
	rows, err := q.db.QueryContext(ctx, getPostImagesUrl, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []sql.NullString
	for rows.Next() {
		var url sql.NullString
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		items = append(items, url)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertLikedPost = `-- name: InsertLikedPost :one
INSERT INTO liked_posts (user_id, post_id)
VALUES ($1::bigint, $2::bigint)
ON CONFLICT (user_id, post_id) DO NOTHING
RETURNING id
`

type InsertLikedPostParams struct {
	UserID int64
	PostID int64
}

func (q *Queries) InsertLikedPost(ctx context.Context, arg InsertLikedPostParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, insertLikedPost, arg.UserID, arg.PostID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const insertPost = `-- name: InsertPost :one
INSERT INTO posts
(user_id, title, content, visibility)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, content, like_count, comment_count, repost_count, created_at, updated_at, title, visibility
`

type InsertPostParams struct {
	UserID     sql.NullInt64
	Title      string
	Content    sql.NullString
	Visibility string
}

func (q *Queries) InsertPost(ctx context.Context, arg InsertPostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, insertPost,
		arg.UserID,
		arg.Title,
		arg.Content,
		arg.Visibility,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Content,
		&i.LikeCount,
		&i.CommentCount,
		&i.RepostCount,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Visibility,
	)
	return i, err
}

const insertReportedPost = `-- name: InsertReportedPost :one
INSERT INTO reported_posts
(user_id, post_id, reason, message)
SELECT $1::bigint, $2::bigint, UNNEST($3::varchar(15)[]), $4::text
RETURNING id, user_id, post_id, reason, message
`

type InsertReportedPostParams struct {
	UserID  int64
	PostID  int64
	Reason  []string
	Message string
}

func (q *Queries) InsertReportedPost(ctx context.Context, arg InsertReportedPostParams) (ReportedPost, error) {
	row := q.db.QueryRowContext(ctx, insertReportedPost,
		arg.UserID,
		arg.PostID,
		pq.Array(arg.Reason),
		arg.Message,
	)
	var i ReportedPost
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PostID,
		&i.Reason,
		&i.Message,
	)
	return i, err
}

const insertRepostedPost = `-- name: InsertRepostedPost :one
INSERT INTO reposted_posts (user_id, post_id)
VALUES ($1::bigint, $2::bigint)
ON CONFLICT (user_id, post_id) DO NOTHING
RETURNING id
`

type InsertRepostedPostParams struct {
	UserID int64
	PostID int64
}

func (q *Queries) InsertRepostedPost(ctx context.Context, arg InsertRepostedPostParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, insertRepostedPost, arg.UserID, arg.PostID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const listLikedPostsByUserId = `-- name: ListLikedPostsByUserId :many
SELECT p.id, p.user_id, p.content, p.like_count, p.comment_count, p.repost_count, p.created_at, p.updated_at, p.title, p.visibility, 
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
LIMIT $3
`

type ListLikedPostsByUserIdParams struct {
	UserID sql.NullInt64
	Offset int32
	Limit  int32
}

type ListLikedPostsByUserIdRow struct {
	ID           int64
	UserID       sql.NullInt64
	Content      sql.NullString
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
	ImageUrls    interface{}
	TotalRows    int64
	Liked        bool
	Repost       bool
}

func (q *Queries) ListLikedPostsByUserId(ctx context.Context, arg ListLikedPostsByUserIdParams) ([]ListLikedPostsByUserIdRow, error) {
	rows, err := q.db.QueryContext(ctx, listLikedPostsByUserId, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListLikedPostsByUserIdRow
	for rows.Next() {
		var i ListLikedPostsByUserIdRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Content,
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
			&i.ImageUrls,
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

const listNewestPostsByUserId = `-- name: ListNewestPostsByUserId :many
SELECT p.id, p.user_id, p.content, p.like_count, p.comment_count, p.repost_count, p.created_at, p.updated_at, p.title, p.visibility, 
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
LIMIT $3
`

type ListNewestPostsByUserIdParams struct {
	UserID sql.NullInt64
	Offset int32
	Limit  int32
}

type ListNewestPostsByUserIdRow struct {
	ID           int64
	UserID       sql.NullInt64
	Content      sql.NullString
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
	ImageUrls    interface{}
	TotalRows    int64
	Liked        bool
	Repost       bool
}

func (q *Queries) ListNewestPostsByUserId(ctx context.Context, arg ListNewestPostsByUserIdParams) ([]ListNewestPostsByUserIdRow, error) {
	rows, err := q.db.QueryContext(ctx, listNewestPostsByUserId, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListNewestPostsByUserIdRow
	for rows.Next() {
		var i ListNewestPostsByUserIdRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Content,
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
			&i.ImageUrls,
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

const listRepostedPostsByUserId = `-- name: ListRepostedPostsByUserId :many
SELECT p.id, p.user_id, p.content, p.like_count, p.comment_count, p.repost_count, p.created_at, p.updated_at, p.title, p.visibility, 
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
LIMIT $3
`

type ListRepostedPostsByUserIdParams struct {
	UserID sql.NullInt64
	Offset int32
	Limit  int32
}

type ListRepostedPostsByUserIdRow struct {
	ID           int64
	UserID       sql.NullInt64
	Content      sql.NullString
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
	ImageUrls    interface{}
	TotalRows    int64
	Liked        bool
	Repost       bool
}

func (q *Queries) ListRepostedPostsByUserId(ctx context.Context, arg ListRepostedPostsByUserIdParams) ([]ListRepostedPostsByUserIdRow, error) {
	rows, err := q.db.QueryContext(ctx, listRepostedPostsByUserId, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListRepostedPostsByUserIdRow
	for rows.Next() {
		var i ListRepostedPostsByUserIdRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Content,
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
			&i.ImageUrls,
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

const lockPostForUpdate = `-- name: LockPostForUpdate :one
SELECT 1
FROM posts
WHERE id = $1
FOR UPDATE
`

func (q *Queries) LockPostForUpdate(ctx context.Context, id int64) (int32, error) {
	row := q.db.QueryRowContext(ctx, lockPostForUpdate, id)
	var column_1 int32
	err := row.Scan(&column_1)
	return column_1, err
}

const updatePost = `-- name: UpdatePost :exec
UPDATE posts
SET title = $1::text,
    content = $2::text,
    visibility = $3::varchar(10)
WHERE id = $4::bigint AND user_id = $5::bigint
`

type UpdatePostParams struct {
	Title      string
	Content    string
	Visibility string
	ID         int64
	UserID     int64
}

func (q *Queries) UpdatePost(ctx context.Context, arg UpdatePostParams) error {
	_, err := q.db.ExecContext(ctx, updatePost,
		arg.Title,
		arg.Content,
		arg.Visibility,
		arg.ID,
		arg.UserID,
	)
	return err
}

const updatePostLikeCount = `-- name: UpdatePostLikeCount :one
UPDATE posts
SET like_count = GREATEST(like_count + $1::smallint, 0)
WHERE id = $2::bigint
RETURNING id, like_count
`

type UpdatePostLikeCountParams struct {
	Value int16
	ID    int64
}

type UpdatePostLikeCountRow struct {
	ID        int64
	LikeCount sql.NullInt32
}

func (q *Queries) UpdatePostLikeCount(ctx context.Context, arg UpdatePostLikeCountParams) (UpdatePostLikeCountRow, error) {
	row := q.db.QueryRowContext(ctx, updatePostLikeCount, arg.Value, arg.ID)
	var i UpdatePostLikeCountRow
	err := row.Scan(&i.ID, &i.LikeCount)
	return i, err
}

const updatePostRepostCount = `-- name: UpdatePostRepostCount :one
UPDATE posts
SET repost_count = GREATEST(repost_count + $1::smallint, 0)
WHERE id = $2::bigint
RETURNING id, repost_count
`

type UpdatePostRepostCountParams struct {
	Value int16
	ID    int64
}

type UpdatePostRepostCountRow struct {
	ID          int64
	RepostCount sql.NullInt32
}

func (q *Queries) UpdatePostRepostCount(ctx context.Context, arg UpdatePostRepostCountParams) (UpdatePostRepostCountRow, error) {
	row := q.db.QueryRowContext(ctx, updatePostRepostCount, arg.Value, arg.ID)
	var i UpdatePostRepostCountRow
	err := row.Scan(&i.ID, &i.RepostCount)
	return i, err
}