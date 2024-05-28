package homepage

import (
	"context"
	"database/sql"
	"profiln-be/model"
	homepageSqlc "profiln-be/package/homepage/repository/sqlc"
)

type IHomepageRepository interface {
	ListNewestPosts(userId int64, offset, limit int32) ([]model.Post, int64, error)
	ListPostsByFollowing(userId int64, offset, limit int32) ([]model.Post, int64, error)
	ListPopularPosts(userId int64, offset, limit int32) ([]model.Post, int64, error)
	GetFollowsRecommendationForUserId(userId int64, offset, limit int32) ([]homepageSqlc.GetFollowsRecommendationForUserIdRow, int64, error)
}

type HomepageRepository struct {
	db    *sql.DB
	query *homepageSqlc.Queries
}

func NewHomepageRepository(db *sql.DB) IHomepageRepository {
	return &HomepageRepository{
		db:    db,
		query: homepageSqlc.New(db),
	}
}

func (r *HomepageRepository) ListNewestPosts(userId int64, offset, limit int32) ([]model.Post, int64, error) {
	arg := homepageSqlc.ListNewestPostsParams{
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
		posts[i] = model.Post{
			ID: v.ID,
			User: model.User{
				ID:         v.UserID.Int64,
				AvatarUrl:  v.AvatarUrl.String,
				Fullname:   v.FullName.String,
				Bio:        v.Bio.String,
				OpenToWork: v.OpenToWork.Bool,
			},
			Title:          v.Title,
			Content:        v.Content.String,
			ImageUrl:       v.ImageUrl.String,
			LikeCount:      v.LikeCount.Int32,
			CommentCount:   v.CommentCount.Int32,
			RepostCount:    v.RepostCount.Int32,
			IsRepost:       v.Repost.Bool,
			IsLiked:        v.Liked,
			OriginalPostID: v.OriginalPostID.Int64,
			UpdatedAt:      v.UpdatedAt.Time,
		}
	}

	return posts, count, nil
}

func (r *HomepageRepository) ListPostsByFollowing(userId int64, offset, limit int32) ([]model.Post, int64, error) {
	arg := homepageSqlc.ListPostsByFollowingParams{
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
		posts[i] = model.Post{
			ID: v.ID,
			User: model.User{
				ID:         v.UserID.Int64,
				AvatarUrl:  v.AvatarUrl.String,
				Fullname:   v.FullName.String,
				Bio:        v.Bio.String,
				OpenToWork: v.OpenToWork.Bool,
			},
			Title:          v.Title,
			Content:        v.Content.String,
			ImageUrl:       v.ImageUrl.String,
			LikeCount:      v.LikeCount.Int32,
			CommentCount:   v.CommentCount.Int32,
			RepostCount:    v.RepostCount.Int32,
			IsRepost:       v.Repost.Bool,
			IsLiked:        v.Liked,
			OriginalPostID: v.OriginalPostID.Int64,
			UpdatedAt:      v.UpdatedAt.Time,
		}
	}

	return posts, count, nil
}

func (r *HomepageRepository) ListPopularPosts(userId int64, offset, limit int32) ([]model.Post, int64, error) {
	arg := homepageSqlc.ListPopularPostsParams{
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
		posts[i] = model.Post{
			ID: v.ID,
			User: model.User{
				ID:         v.UserID.Int64,
				AvatarUrl:  v.AvatarUrl.String,
				Fullname:   v.FullName.String,
				Bio:        v.Bio.String,
				OpenToWork: v.OpenToWork.Bool,
			},
			Title:          v.Title,
			Content:        v.Content.String,
			ImageUrl:       v.ImageUrl.String,
			LikeCount:      v.LikeCount.Int32,
			CommentCount:   v.CommentCount.Int32,
			RepostCount:    v.RepostCount.Int32,
			IsRepost:       v.Repost.Bool,
			IsLiked:        v.Liked,
			OriginalPostID: v.OriginalPostID.Int64,
			UpdatedAt:      v.UpdatedAt.Time,
		}
	}

	return posts, count, nil
}

func (r *HomepageRepository) GetFollowsRecommendationForUserId(userId int64, offset, limit int32) ([]homepageSqlc.GetFollowsRecommendationForUserIdRow, int64, error) {
	arg := homepageSqlc.GetFollowsRecommendationForUserIdParams{
		UserID: sql.NullInt64{Int64: userId, Valid: true},
		Offset: offset,
		Limit:  limit,
	}

	data, err := r.query.GetFollowsRecommendationForUserId(context.Background(), arg)
	if err != nil {
		return []homepageSqlc.GetFollowsRecommendationForUserIdRow{}, 0, err
	}

	// get total rows for pagination
	var count int64
	if len(data) > 0 {
		count = data[0].TotalRows
	}

	return data, count, nil
}
