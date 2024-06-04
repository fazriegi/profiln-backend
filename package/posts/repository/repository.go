package posts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	db "profiln-be/db/sqlc"
	"profiln-be/model"

	"strings"
	"sync"
)

type IPostsRepository interface {
	InsertReportedPost(userId int64, props *model.ReportPost) (db.ReportedPost, error)
	GetDetailPost(postId, userId int64) (model.Post, error)
	GetPostComments(postId int64, offset, limit int32) ([]db.GetPostCommentsRow, int64, error)
	GetPostCommentReplies(postId, postCommentId int64, offset, limit int32) ([]db.GetPostCommentRepliesRow, int64, error)
	LikePost(userId, postId int64) (*db.UpdatePostLikeCountRow, error)
	UnlikePost(userId, postId int64) (*db.UpdatePostLikeCountRow, error)
	ListNewestPostsByUserId(userId int64, offset, limit int32) ([]model.Post, int64, error)
	ListLikedPostsByUserId(userId int64, offset, limit int32) ([]model.Post, int64, error)
	ListRepostedPostsByUserId(userId int64, offset, limit int32) ([]model.Post, int64, error)
	InsertPost(props *model.CreatePostRequest) (model.Post, error)
	UpdatePostById(props *model.UpdatePostRequest) error
	GetPostById(postId int64) (model.Post, error)
	GetPostImagesUrl(postId int64) ([]string, error)
	DeletePost(postId int64) error
	RepostPost(userId, postId int64) (*db.UpdatePostRepostCountRow, error)
	UnrepostPost(userId, postId int64) (*db.UpdatePostRepostCountRow, error)
	BatchInsertPostImages(postId int64, urls []string) ([]db.PostImage, error)
	CountPostImages(postId int64) (int64, error)
}

type PostsRepository struct {
	dbConn *sql.DB
	query  *db.Queries
}

func NewPostsRepository(dbConn *sql.DB) IPostsRepository {
	return &PostsRepository{
		dbConn: dbConn,
		query:  db.New(dbConn),
	}
}

func (r *PostsRepository) InsertReportedPost(userId int64, props *model.ReportPost) (db.ReportedPost, error) {
	arg := db.InsertReportedPostParams{
		UserID:  userId,
		PostID:  props.PostId,
		Reason:  props.Reason,
		Message: props.Message,
	}

	reportedPost, err := r.query.InsertReportedPost(context.Background(), arg)
	if err != nil {
		return db.ReportedPost{}, err
	}

	return reportedPost, nil
}

func (r *PostsRepository) GetDetailPost(postId, userId int64) (model.Post, error) {
	data, err := r.query.GetDetailPost(context.Background(), db.GetDetailPostParams{
		ID:     postId,
		UserID: sql.NullInt64{Int64: userId, Valid: true},
	})
	if err != nil {
		return model.Post{}, err
	}

	var imageUrls []string

	// Convert to array
	if data.ImageUrls != nil {
		imageUrlsString := strings.Trim(string(data.ImageUrls.([]uint8)), "{}")
		imageUrls = strings.Split(imageUrlsString, ",")
	}

	post := model.Post{
		ID: data.ID,
		User: model.User{
			ID:         data.ID_2,
			Fullname:   data.FullName,
			AvatarUrl:  data.AvatarUrl.String,
			Bio:        data.Bio.String,
			OpenToWork: data.OpenToWork.Bool,
		},
		Title:        data.Title,
		Content:      data.Content.String,
		ImageUrls:    imageUrls,
		LikeCount:    data.LikeCount.Int32,
		CommentCount: data.CommentCount.Int32,
		RepostCount:  data.RepostCount.Int32,
		IsRepost:     data.Repost,
		IsLiked:      data.Liked,
		UpdatedAt:    data.UpdatedAt.Time,
	}

	return post, nil
}

func (r *PostsRepository) GetPostComments(postId int64, offset, limit int32) ([]db.GetPostCommentsRow, int64, error) {
	arg := db.GetPostCommentsParams{
		PostID: sql.NullInt64{Int64: postId, Valid: true},
		Offset: offset,
		Limit:  limit,
	}

	data, err := r.query.GetPostComments(context.Background(), arg)
	if err != nil {
		return []db.GetPostCommentsRow{}, 0, err
	}

	var count int64
	if len(data) > 0 {
		count = data[0].TotalRows
	}

	return data, count, nil
}

func (r *PostsRepository) GetPostCommentReplies(postId, postCommentId int64, offset, limit int32) ([]db.GetPostCommentRepliesRow, int64, error) {
	arg := db.GetPostCommentRepliesParams{
		PostID:        sql.NullInt64{Int64: postId, Valid: true},
		PostCommentID: sql.NullInt64{Int64: postCommentId, Valid: true},
		Offset:        offset,
		Limit:         limit,
	}

	data, err := r.query.GetPostCommentReplies(context.Background(), arg)
	if err != nil {
		return []db.GetPostCommentRepliesRow{}, 0, err
	}

	var count int64
	if len(data) > 0 {
		count = data[0].TotalRows
	}

	return data, count, nil
}

func (r *PostsRepository) LikePost(userId, postId int64) (*db.UpdatePostLikeCountRow, error) {
	ctx := context.Background()
	tx, err := r.dbConn.Begin()
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

	post, err := qtx.UpdatePostLikeCount(ctx, db.UpdatePostLikeCountParams{
		ID:    postId,
		Value: 1,
	})
	if err != nil {
		return nil, fmt.Errorf("could not update like count: %w", err)
	}

	_, err = qtx.InsertLikedPost(ctx, db.InsertLikedPostParams{
		UserID: userId,
		PostID: postId,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			post = db.UpdatePostLikeCountRow{
				ID:        post.ID,
				LikeCount: sql.NullInt32{Int32: post.LikeCount.Int32 - 1, Valid: true},
			}

			return &post, nil
		}

		return nil, fmt.Errorf("could not insert liked post: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	return &post, nil
}

func (r *PostsRepository) UnlikePost(userId, postId int64) (*db.UpdatePostLikeCountRow, error) {
	ctx := context.Background()
	tx, err := r.dbConn.Begin()
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

	post, err := qtx.UpdatePostLikeCount(ctx, db.UpdatePostLikeCountParams{
		ID:    postId,
		Value: -1,
	})
	if err != nil {
		return nil, fmt.Errorf("could not update like count: %w", err)
	}

	_, err = qtx.DeleteLikedPost(ctx, db.DeleteLikedPostParams{
		UserID: userId,
		PostID: postId,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			post = db.UpdatePostLikeCountRow{
				ID:        post.ID,
				LikeCount: sql.NullInt32{Int32: post.LikeCount.Int32 + 1, Valid: true},
			}

			return &post, nil
		}

		return nil, fmt.Errorf("could not delete liked post: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	return &post, nil
}

func (r *PostsRepository) ListNewestPostsByUserId(userId int64, offset, limit int32) ([]model.Post, int64, error) {
	arg := db.ListNewestPostsByUserIdParams{
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

func (r *PostsRepository) ListLikedPostsByUserId(userId int64, offset, limit int32) ([]model.Post, int64, error) {
	arg := db.ListLikedPostsByUserIdParams{
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

func (r *PostsRepository) ListRepostedPostsByUserId(userId int64, offset, limit int32) ([]model.Post, int64, error) {
	arg := db.ListRepostedPostsByUserIdParams{
		UserID: sql.NullInt64{Int64: userId, Valid: true},
		Offset: offset,
		Limit:  limit,
	}

	data, err := r.query.ListRepostedPostsByUserId(context.Background(), arg)
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

func (r *PostsRepository) InsertPost(props *model.CreatePostRequest) (model.Post, error) {
	insertPostArg := db.InsertPostParams{
		UserID:     sql.NullInt64{Int64: props.UserId, Valid: true},
		Title:      props.Title,
		Content:    sql.NullString{String: props.Content, Valid: true},
		Visibility: props.Visibility,
	}

	createdPost, err := r.query.InsertPost(context.Background(), insertPostArg)
	if err != nil {
		return model.Post{}, err
	}

	data := model.Post{
		ID: createdPost.ID,
		User: model.User{
			ID: createdPost.UserID.Int64,
		},
		Title:        createdPost.Title,
		Content:      createdPost.Content.String,
		LikeCount:    createdPost.LikeCount.Int32,
		CommentCount: createdPost.CommentCount.Int32,
		RepostCount:  createdPost.RepostCount.Int32,
		IsRepost:     false,
		IsLiked:      false,
		UpdatedAt:    createdPost.UpdatedAt.Time,
	}

	return data, nil
}

func (r *PostsRepository) UpdatePostById(props *model.UpdatePostRequest) error {
	// ctx := context.Background()
	// tx, err := r.dbConn.Begin()
	// if err != nil {
	// 	return fmt.Errorf("could not begin transaction: %w", err)
	// }
	// defer tx.Rollback()

	// qtx := r.query.WithTx(tx)

	// err = qtx.BatchDeletePostImagesByPost(ctx, props.ID)
	// if err != nil {
	// 	return fmt.Errorf("could not batch delete post images: %w", err)
	// }

	updatePostArg := db.UpdatePostParams{
		ID:         props.ID,
		UserID:     props.UserId,
		Title:      props.Title,
		Content:    props.Content,
		Visibility: props.Visibility,
	}
	if err := r.query.UpdatePost(context.Background(), updatePostArg); err != nil {
		return fmt.Errorf("could not update post: %w", err)
	}

	// _, err = qtx.BatchInsertPostImages(ctx, db.BatchInsertPostImagesParams{
	// 	PostID: props.ID,
	// 	Url:    props.ImageUrls,
	// })
	// if err != nil {
	// 	return fmt.Errorf("could not insert post images: %w", err)
	// }

	// if err := tx.Commit(); err != nil {
	// 	return fmt.Errorf("could not commit transaction: %w", err)
	// }

	return nil
}

func (r *PostsRepository) GetPostById(postId int64) (model.Post, error) {
	data, err := r.query.GetPostById(context.Background(), postId)
	if err != nil {
		return model.Post{}, err
	}

	post := model.Post{
		ID: data.ID,
		User: model.User{
			ID: data.UserID.Int64,
		},
		Title:        data.Title,
		Content:      data.Content.String,
		LikeCount:    data.LikeCount.Int32,
		CommentCount: data.CommentCount.Int32,
		RepostCount:  data.RepostCount.Int32,
		UpdatedAt:    data.UpdatedAt.Time,
	}

	return post, nil
}

func (r *PostsRepository) GetPostImagesUrl(postId int64) ([]string, error) {
	data, err := r.query.GetPostImagesUrl(context.Background(), postId)
	if err != nil {
		return nil, err
	}

	urls := make([]string, len(data))
	for i, v := range data {
		urls[i] = v.String
	}

	return urls, nil
}

func (r *PostsRepository) DeletePost(postId int64) error {
	var (
		errChan = make(chan error, 5)
		wg      sync.WaitGroup
	)

	ctx := context.Background()
	tx, err := r.dbConn.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.query.WithTx(tx)

	deleteFuncs := []func(int64){
		func(postId int64) {
			defer wg.Done()
			if err := qtx.BatchDeleteReportedPostsByPost(ctx, postId); err != nil {
				errChan <- fmt.Errorf("could not batch delete reported post: %w", err)
			}
		},
		func(postId int64) {
			defer wg.Done()
			if err := qtx.BatchDeleteLikedPostByPost(ctx, postId); err != nil {
				errChan <- fmt.Errorf("could not batch delete liked post: %w", err)
			}
		},
		func(postId int64) {
			defer wg.Done()
			if err := qtx.BatchDeleteRepostedPostByPost(ctx, postId); err != nil {
				errChan <- fmt.Errorf("could not batch delete reposted post: %w", err)
			}
		},
		func(postId int64) {
			defer wg.Done()
			if err := qtx.BatchDeletePostCommentRepliesByPost(ctx, postId); err != nil {
				errChan <- fmt.Errorf("could not batch delete post comment replies: %w", err)
			}
		},
		func(postId int64) {
			defer wg.Done()
			if err := qtx.BatchDeletePostImagesByPost(ctx, postId); err != nil {
				errChan <- fmt.Errorf("could not batch delete post images: %w", err)
			}
		},
	}

	for _, deleteFunc := range deleteFuncs {
		wg.Add(1)
		go deleteFunc(postId)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	if err = qtx.BatchDeletePostCommentsByPost(ctx, postId); err != nil {
		return fmt.Errorf("could not batch delete post comments: %w", err)
	}

	if err = qtx.DeletePostById(ctx, postId); err != nil {
		return fmt.Errorf("could not delete post: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (r *PostsRepository) RepostPost(userId, postId int64) (*db.UpdatePostRepostCountRow, error) {
	ctx := context.Background()
	tx, err := r.dbConn.Begin()
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.query.WithTx(tx)

	_, err = qtx.LockPostForUpdate(ctx, postId)
	if err != nil {
		return nil, fmt.Errorf("could not lock post for update: %w", err)
	}

	post, err := qtx.UpdatePostRepostCount(ctx, db.UpdatePostRepostCountParams{
		ID:    postId,
		Value: 1,
	})
	if err != nil {
		return nil, fmt.Errorf("could not update repost count: %w", err)
	}

	_, err = qtx.InsertRepostedPost(ctx, db.InsertRepostedPostParams{
		UserID: userId,
		PostID: postId,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			post = db.UpdatePostRepostCountRow{
				ID:          post.ID,
				RepostCount: sql.NullInt32{Int32: post.RepostCount.Int32 - 1, Valid: true},
			}

			return &post, nil
		}

		return nil, fmt.Errorf("could not insert reposted post: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	return &post, nil
}

func (r *PostsRepository) UnrepostPost(userId, postId int64) (*db.UpdatePostRepostCountRow, error) {
	ctx := context.Background()
	tx, err := r.dbConn.Begin()
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

	post, err := qtx.UpdatePostRepostCount(ctx, db.UpdatePostRepostCountParams{
		ID:    postId,
		Value: -1,
	})
	if err != nil {
		return nil, fmt.Errorf("could not update repost count: %w", err)
	}

	_, err = qtx.DeleteRepostedPost(ctx, db.DeleteRepostedPostParams{
		UserID: userId,
		PostID: postId,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			post = db.UpdatePostRepostCountRow{
				ID:          post.ID,
				RepostCount: sql.NullInt32{Int32: post.RepostCount.Int32 + 1, Valid: true},
			}

			return &post, nil
		}

		return nil, fmt.Errorf("could not delete reposted post: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	return &post, nil
}

func (r *PostsRepository) BatchInsertPostImages(postId int64, urls []string) ([]db.PostImage, error) {
	ctx := context.Background()
	tx, err := r.dbConn.Begin()
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.query.WithTx(tx)

	err = qtx.BatchDeletePostImagesByPost(ctx, postId)
	if err != nil {
		return nil, fmt.Errorf("could not batch delete post images: %w", err)
	}

	urlIndex := []int16{}
	for i := range urls {
		urlIndex = append(urlIndex, int16(i))
	}

	data, err := qtx.BatchInsertPostImages(ctx, db.BatchInsertPostImagesParams{
		PostID: postId,
		Index:  urlIndex,
		Url:    urls,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	return data, nil
}

func (r *PostsRepository) CountPostImages(postId int64) (int64, error) {
	count, err := r.query.CountPostImages(context.Background(), postId)
	if err != nil {
		return 0, err
	}

	return count, nil
}
