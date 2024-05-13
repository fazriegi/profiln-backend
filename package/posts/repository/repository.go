package posts

import (
	"context"
	"database/sql"
	postSqlc "profiln-be/package/posts/repository/sqlc"
)

type IPostsRepository interface {
	InsertReportedPost(userId, postId int64, reason, message string) (postSqlc.ReportedPost, error)
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
