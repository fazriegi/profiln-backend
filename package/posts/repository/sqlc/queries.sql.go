// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package posts

import (
	"context"
	"database/sql"
)

const getDetailPost = `-- name: GetDetailPost :one
SELECT p.id, p.user_id, p.content, p.image_url, p.like_count, p.comment_count, p.repost_count, p.repost, p.original_post_id, p.created_at, p.updated_at, 
    pu.id, pu.avatar_url, pu.full_name, pu.bio, pu.open_to_work
FROM posts p
JOIN users pu ON p.user_id = pu.id
WHERE p.id = $1
`

type GetDetailPostRow struct {
	ID             int64
	UserID         sql.NullInt64
	Content        sql.NullString
	ImageUrl       sql.NullString
	LikeCount      sql.NullInt32
	CommentCount   sql.NullInt32
	RepostCount    sql.NullInt32
	Repost         sql.NullBool
	OriginalPostID sql.NullInt64
	CreatedAt      sql.NullTime
	UpdatedAt      sql.NullTime
	ID_2           int64
	AvatarUrl      sql.NullString
	FullName       string
	Bio            sql.NullString
	OpenToWork     sql.NullBool
}

func (q *Queries) GetDetailPost(ctx context.Context, id int64) (GetDetailPostRow, error) {
	row := q.db.QueryRowContext(ctx, getDetailPost, id)
	var i GetDetailPostRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Content,
		&i.ImageUrl,
		&i.LikeCount,
		&i.CommentCount,
		&i.RepostCount,
		&i.Repost,
		&i.OriginalPostID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ID_2,
		&i.AvatarUrl,
		&i.FullName,
		&i.Bio,
		&i.OpenToWork,
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

const insertReportedPost = `-- name: InsertReportedPost :one
INSERT INTO reported_posts
(user_id, post_id, reason, message)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, post_id, reason, message
`

type InsertReportedPostParams struct {
	UserID  sql.NullInt64
	PostID  sql.NullInt64
	Reason  sql.NullString
	Message sql.NullString
}

func (q *Queries) InsertReportedPost(ctx context.Context, arg InsertReportedPostParams) (ReportedPost, error) {
	row := q.db.QueryRowContext(ctx, insertReportedPost,
		arg.UserID,
		arg.PostID,
		arg.Reason,
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

const updatePostLikeCount = `-- name: UpdatePostLikeCount :one
UPDATE posts
SET like_count = like_count + 1
WHERE id = $1
RETURNING id, like_count
`

type UpdatePostLikeCountRow struct {
	ID        int64
	LikeCount sql.NullInt32
}

func (q *Queries) UpdatePostLikeCount(ctx context.Context, id int64) (UpdatePostLikeCountRow, error) {
	row := q.db.QueryRowContext(ctx, updatePostLikeCount, id)
	var i UpdatePostLikeCountRow
	err := row.Scan(&i.ID, &i.LikeCount)
	return i, err
}
