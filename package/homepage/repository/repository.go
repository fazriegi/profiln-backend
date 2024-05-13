package homepage

import (
	"context"
	"database/sql"
	homepageSqlc "profiln-be/package/homepage/repository/sqlc"
)

type IHomepageRepository interface {
	ListPosts(userId int64, offset, limit int32) ([]homepageSqlc.Post, int64, error)
	ListPostsByFollowing(userId int64, offset, limit int32) ([]homepageSqlc.Post, int64, error)
	ListPopularPosts(userId int64, offset, limit int32) ([]homepageSqlc.Post, int64, error)
	GetUserById(id int64) (homepageSqlc.User, error)
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

func (r *HomepageRepository) ListPosts(userId int64, offset, limit int32) ([]homepageSqlc.Post, int64, error) {
	posts := []homepageSqlc.Post{}
	arg := homepageSqlc.ListPostsParams{
		UserID: sql.NullInt64{Int64: userId, Valid: true},
		Offset: offset,
		Limit:  limit,
	}

	data, err := r.query.ListPosts(context.Background(), arg)
	if err != nil {
		return []homepageSqlc.Post{}, 0, err
	}

	var count int64
	if len(data) > 0 {
		count = data[0].TotalRows
	}

	for _, v := range data {
		post := homepageSqlc.Post{}
		post.ID = v.ID
		post.UserID = v.UserID
		post.Content = v.Content
		post.ImageUrl = v.ImageUrl
		post.LikeCount = v.LikeCount
		post.CommentCount = v.CommentCount
		post.RepostCount = v.RepostCount
		post.Repost = v.Repost
		post.OriginalPostID = v.OriginalPostID
		post.CreatedAt = v.CreatedAt
		post.UpdatedAt = v.UpdatedAt

		posts = append(posts, post)
	}

	return posts, count, nil
}

func (r *HomepageRepository) ListPostsByFollowing(userId int64, offset, limit int32) ([]homepageSqlc.Post, int64, error) {
	posts := []homepageSqlc.Post{}
	arg := homepageSqlc.ListPostsByFollowingParams{
		UserID: sql.NullInt64{Int64: userId, Valid: true},
		Offset: offset,
		Limit:  limit,
	}

	data, err := r.query.ListPostsByFollowing(context.Background(), arg)
	if err != nil {
		return []homepageSqlc.Post{}, 0, err
	}

	var count int64
	if len(data) > 0 {
		count = data[0].TotalRows
	}

	for _, v := range data {
		post := homepageSqlc.Post{}
		post.ID = v.ID
		post.UserID = v.UserID
		post.Content = v.Content
		post.ImageUrl = v.ImageUrl
		post.LikeCount = v.LikeCount
		post.CommentCount = v.CommentCount
		post.RepostCount = v.RepostCount
		post.Repost = v.Repost
		post.OriginalPostID = v.OriginalPostID
		post.CreatedAt = v.CreatedAt
		post.UpdatedAt = v.UpdatedAt

		posts = append(posts, post)
	}

	return posts, count, nil
}

func (r *HomepageRepository) ListPopularPosts(userId int64, offset, limit int32) ([]homepageSqlc.Post, int64, error) {
	posts := []homepageSqlc.Post{}
	arg := homepageSqlc.ListPopularPostsParams{
		UserID: sql.NullInt64{Int64: userId, Valid: true},
		Offset: offset,
		Limit:  limit,
	}

	data, err := r.query.ListPopularPosts(context.Background(), arg)
	if err != nil {
		return []homepageSqlc.Post{}, 0, err
	}

	var count int64
	if len(data) > 0 {
		count = data[0].TotalRows
	}

	for _, v := range data {
		post := homepageSqlc.Post{}
		post.ID = v.ID
		post.UserID = v.UserID
		post.Content = v.Content
		post.ImageUrl = v.ImageUrl
		post.LikeCount = v.LikeCount
		post.CommentCount = v.CommentCount
		post.RepostCount = v.RepostCount
		post.Repost = v.Repost
		post.OriginalPostID = v.OriginalPostID
		post.CreatedAt = v.CreatedAt
		post.UpdatedAt = v.UpdatedAt

		posts = append(posts, post)
	}

	return posts, count, nil
}

func (r *HomepageRepository) GetUserById(id int64) (homepageSqlc.User, error) {
	user, err := r.query.GetUserById(context.Background(), id)
	if err != nil {
		return homepageSqlc.User{}, err
	}

	return user, nil
}
