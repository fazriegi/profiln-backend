package homepage

import (
	"context"
	"database/sql"
	db "profiln-be/db/sqlc"
	"profiln-be/model"

	"strings"
)

type IHomepageRepository interface {
	ListNewestPosts(userId int64, offset, limit int32) ([]model.Post, int64, error)
	ListPostsByFollowing(userId int64, offset, limit int32) ([]model.Post, int64, error)
	ListPopularPosts(userId int64, offset, limit int32) ([]model.Post, int64, error)
	GetFollowsRecommendationForUserId(userId int64, offset, limit int32) ([]db.GetFollowsRecommendationForUserIdRow, int64, error)
}

type HomepageRepository struct {
	dbConn *sql.DB
	query  *db.Queries
}

func NewHomepageRepository(dbConn *sql.DB) IHomepageRepository {
	return &HomepageRepository{
		dbConn: dbConn,
		query:  db.New(dbConn),
	}
}

func (r *HomepageRepository) ListNewestPosts(userId int64, offset, limit int32) ([]model.Post, int64, error) {
	arg := db.ListNewestPostsParams{
		UserID: sql.NullInt64{Int64: userId, Valid: true},
		Offset: offset,
		Limit:  limit,
	}

	data, err := r.query.ListNewestPosts(context.Background(), arg)
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
		var imageUrls []string

		// Convert to array
		if v.ImageUrls != nil {
			imageUrlsString := strings.Trim(string(v.ImageUrls.([]uint8)), "{}")
			imageUrls = strings.Split(imageUrlsString, ",")
		}

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
			ImageUrls:    imageUrls,
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

func (r *HomepageRepository) ListPostsByFollowing(userId int64, offset, limit int32) ([]model.Post, int64, error) {
	arg := db.ListPostsByFollowingParams{
		UserID: sql.NullInt64{Int64: userId, Valid: true},
		Offset: offset,
		Limit:  limit,
	}

	data, err := r.query.ListPostsByFollowing(context.Background(), arg)
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
		var imageUrls []string

		// Convert to array
		if v.ImageUrls != nil {
			imageUrlsString := strings.Trim(string(v.ImageUrls.([]uint8)), "{}")
			imageUrls = strings.Split(imageUrlsString, ",")
		}

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
			ImageUrls:    imageUrls,
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

func (r *HomepageRepository) ListPopularPosts(userId int64, offset, limit int32) ([]model.Post, int64, error) {
	arg := db.ListPopularPostsParams{
		UserID: sql.NullInt64{Int64: userId, Valid: true},
		Offset: offset,
		Limit:  limit,
	}

	data, err := r.query.ListPopularPosts(context.Background(), arg)
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
		var imageUrls []string

		// Convert to array
		if v.ImageUrls != nil {
			imageUrlsString := strings.Trim(string(v.ImageUrls.([]uint8)), "{}")
			imageUrls = strings.Split(imageUrlsString, ",")
		}

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
			ImageUrls:    imageUrls,
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

func (r *HomepageRepository) GetFollowsRecommendationForUserId(userId int64, offset, limit int32) ([]db.GetFollowsRecommendationForUserIdRow, int64, error) {
	arg := db.GetFollowsRecommendationForUserIdParams{
		UserID: sql.NullInt64{Int64: userId, Valid: true},
		Offset: offset,
		Limit:  limit,
	}

	data, err := r.query.GetFollowsRecommendationForUserId(context.Background(), arg)
	if err != nil {
		return []db.GetFollowsRecommendationForUserIdRow{}, 0, err
	}

	// get total rows for pagination
	var count int64
	if len(data) > 0 {
		count = data[0].TotalRows
	}

	return data, count, nil
}
