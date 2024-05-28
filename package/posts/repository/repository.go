package posts

import (
	"context"
	"database/sql"
	"fmt"
	"profiln-be/model"
	postSqlc "profiln-be/package/posts/repository/sqlc"
)

type IPostsRepository interface {
	InsertReportedPost(userId, postId int64, reason, message string) (postSqlc.ReportedPost, error)
	GetDetailPost(postId, userId int64) (postSqlc.GetDetailPostRow, error)
	GetPostComments(postId int64, offset, limit int32) ([]postSqlc.GetPostCommentsRow, int64, error)
	GetPostCommentReplies(postId, postCommentId int64, offset, limit int32) ([]postSqlc.GetPostCommentRepliesRow, int64, error)
	UpdatePostLikeCount(postId int64) (*postSqlc.UpdatePostLikeCountRow, error)
	ListNewestPostsByUserId(userId int64, offset, limit int32) ([]model.Post, int64, error)
	ListLikedPostsByUserId(userId int64, offset, limit int32) ([]model.Post, int64, error)
}

type PostsRepository struct {
	db    *sql.DB
	query *postSqlc.Queries
}

func NewPostsRepository(db *sql.DB) IPostsRepository {
	return &PostsRepository{
		db:    db,
		query: postSqlc.New(db),
	}
}

func (r *PostsRepository) InsertReportedPost(userId, postId int64, reason, message string) (postSqlc.ReportedPost, error) {
	arg := postSqlc.InsertReportedPostParams{
		UserID:  sql.NullInt64{Int64: userId, Valid: true},
		PostID:  sql.NullInt64{Int64: postId, Valid: true},
		Reason:  sql.NullString{String: reason, Valid: true},
		Message: sql.NullString{String: message, Valid: true},
	}

	reportedPost, err := r.query.InsertReportedPost(context.Background(), arg)
	if err != nil {
		return postSqlc.ReportedPost{}, err
	}

	return reportedPost, nil
}

func (r *PostsRepository) GetDetailPost(postId, userId int64) (postSqlc.GetDetailPostRow, error) {
	data, err := r.query.GetDetailPost(context.Background(), postSqlc.GetDetailPostParams{
		ID:     postId,
		UserID: sql.NullInt64{Int64: userId, Valid: true},
	})
	if err != nil {
		return postSqlc.GetDetailPostRow{}, err
	}

	return data, nil
}

func (r *PostsRepository) GetPostComments(postId int64, offset, limit int32) ([]postSqlc.GetPostCommentsRow, int64, error) {
	arg := postSqlc.GetPostCommentsParams{
		PostID: sql.NullInt64{Int64: postId, Valid: true},
		Offset: offset,
		Limit:  limit,
	}

	data, err := r.query.GetPostComments(context.Background(), arg)
	if err != nil {
		return []postSqlc.GetPostCommentsRow{}, 0, err
	}

	var count int64
	if len(data) > 0 {
		count = data[0].TotalRows
	}

	return data, count, nil
}

func (r *PostsRepository) GetPostCommentReplies(postId, postCommentId int64, offset, limit int32) ([]postSqlc.GetPostCommentRepliesRow, int64, error) {
	arg := postSqlc.GetPostCommentRepliesParams{
		PostID:        sql.NullInt64{Int64: postId, Valid: true},
		PostCommentID: sql.NullInt64{Int64: postCommentId, Valid: true},
		Offset:        offset,
		Limit:         limit,
	}

	data, err := r.query.GetPostCommentReplies(context.Background(), arg)
	if err != nil {
		return []postSqlc.GetPostCommentRepliesRow{}, 0, err
	}

	var count int64
	if len(data) > 0 {
		count = data[0].TotalRows
	}

	return data, count, nil
}

func (r *PostsRepository) UpdatePostLikeCount(postId int64) (*postSqlc.UpdatePostLikeCountRow, error) {
	ctx := context.Background()
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.query.WithTx(tx)

	_, err = qtx.LockPostForUpdate(ctx, postId)
	if err != nil && err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("could not lock post for update: %w", err)
	}

	post, err := qtx.UpdatePostLikeCount(ctx, postId)
	if err != nil {
		return nil, fmt.Errorf("could not update like count: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	return &post, nil
}

func (r *PostsRepository) ListNewestPostsByUserId(userId int64, offset, limit int32) ([]model.Post, int64, error) {
	arg := postSqlc.ListNewestPostsByUserIdParams{
		UserID: sql.NullInt64{Int64: userId, Valid: true},
		Offset: offset,
		Limit:  limit,
	}
	fmt.Println(arg)

	data, err := r.query.ListNewestPostsByUserId(context.Background(), arg)
	if err != nil {
		return []model.Post{}, 0, err
	}

	// get total rows for pagination
	var count int64
	if len(data) > 0 {
		count = data[0].TotalRows
	}

	posts := make([]model.Post, len(data))
	for i, v := range data {
		posts[i] = model.Post{
			ID: v.ID,
			User: model.User{
				ID:         v.UserID.Int64,
				AvatarUrl:  v.AvatarUrl.String,
				Fullname:   v.FullName.String,
				Bio:        v.Bio.String,
				OpenToWork: v.OpenToWork.Bool,
			},
			Title:        v.Title,
			Content:      v.Content.String,
			ImageUrl:     v.ImageUrl.String,
			LikeCount:    v.LikeCount.Int32,
			CommentCount: v.CommentCount.Int32,
			RepostCount:  v.RepostCount.Int32,
			IsRepost:     v.Repost,
			IsLiked:      v.Liked,
			UpdatedAt:    v.UpdatedAt.Time,
		}
	}

	return posts, count, nil
}

func (r *PostsRepository) ListLikedPostsByUserId(userId int64, offset, limit int32) ([]model.Post, int64, error) {
	arg := postSqlc.ListLikedPostsByUserIdParams{
		UserID: sql.NullInt64{Int64: userId, Valid: true},
		Offset: offset,
		Limit:  limit,
	}

	data, err := r.query.ListLikedPostsByUserId(context.Background(), arg)
	if err != nil {
		return []model.Post{}, 0, err
	}

	// get total rows for pagination
	var count int64
	if len(data) > 0 {
		count = data[0].TotalRows
	}

	posts := make([]model.Post, len(data))
	for i, v := range data {
		posts[i] = model.Post{
			ID: v.ID,
			User: model.User{
				ID:         v.UserID.Int64,
				AvatarUrl:  v.AvatarUrl.String,
				Fullname:   v.FullName.String,
				Bio:        v.Bio.String,
				OpenToWork: v.OpenToWork.Bool,
			},
			Title:        v.Title,
			Content:      v.Content.String,
			ImageUrl:     v.ImageUrl.String,
			LikeCount:    v.LikeCount.Int32,
			CommentCount: v.CommentCount.Int32,
			RepostCount:  v.RepostCount.Int32,
			IsRepost:     v.Repost,
			IsLiked:      v.Liked,
			UpdatedAt:    v.UpdatedAt.Time,
		}
	}

	return posts, count, nil
}
