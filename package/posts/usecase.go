package posts

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"net/http"
	"profiln-be/libs"
	"profiln-be/model"
	repository "profiln-be/package/posts/repository"
	"sync"

	"github.com/sirupsen/logrus"
)

type IPostsUsecase interface {
	InsertReportedPost(userId int64, data *model.ReportPostRequest) (resp model.Response)
	GetDetailPost(postId, userId int64) (resp model.Response)
	GetPostComments(postId int64, pagination model.PaginationRequest) (resp model.Response)
	GetPostCommentReplies(postId, postCommentId int64, pagination model.PaginationRequest) (resp model.Response)
	UpdatePostLikeCount(postId int64) (resp model.Response)
	ListNewestPostsByUserId(userId int64, pagination model.PaginationRequest) (resp model.Response)
	ListLikedPostsByUserId(userId int64, pagination model.PaginationRequest) (resp model.Response)
	ListRepostedPostsByUserId(userId int64, pagination model.PaginationRequest) (resp model.Response)
	InsertPost(imageFile *multipart.FileHeader, props *model.CreatePostRequest) model.Response
}

type PostsUsecase struct {
	repository   repository.IPostsRepository
	log          *logrus.Logger
	googleBucket libs.IGoogleBucket
}

func NewPostsUsecase(repository repository.IPostsRepository, log *logrus.Logger, googleBucket libs.IGoogleBucket) IPostsUsecase {
	return &PostsUsecase{
		repository,
		log,
		googleBucket,
	}
}

func (u *PostsUsecase) InsertReportedPost(userId int64, props *model.ReportPostRequest) (resp model.Response) {
	reportedPost, err :=
		u.repository.InsertReportedPost(userId, props.PostId, props.Reason, props.Message)
	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")
		u.log.Errorf("repository.InsertReportedPost: %v", err)
		return
	}

	data := model.ReportPostResponse{
		PostId:  reportedPost.PostID.Int64,
		Reason:  reportedPost.Reason.String,
		Message: reportedPost.Message.String,
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "Success report post")
	resp.Data = data
	return
}

func (u *PostsUsecase) GetDetailPost(postId, userId int64) (resp model.Response) {
	data, err := u.repository.GetDetailPost(postId, userId)

	if err != nil && err == sql.ErrNoRows {
		return model.Response{
			Status: libs.CustomResponse(http.StatusNotFound, "Data not found"),
		}
	} else if err != nil {
		u.log.Errorf("repository.GetDetailPost: %v", err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occurred"),
		}
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
		LikeCount:    data.LikeCount.Int32,
		CommentCount: data.CommentCount.Int32,
		RepostCount:  data.RepostCount.Int32,
		IsRepost:     data.Repost,
		IsLiked:      data.Liked,
		UpdatedAt:    data.UpdatedAt.Time,
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "Success get detail post")
	resp.Data = post
	return
}

func (u *PostsUsecase) GetPostComments(postId int64, pagination model.PaginationRequest) (resp model.Response) {
	offset := (pagination.Page - 1) * pagination.Limit
	data, totalRows, err := u.repository.GetPostComments(postId, int32(offset), int32(pagination.Limit))

	if err != nil {
		u.log.Errorf("repository.GetPostComments: %v", err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occurred"),
		}
	}

	postComments := make([]model.PostComment, len(data))
	for i, v := range data {
		postComments[i] = model.PostComment{
			ID:     v.ID,
			PostId: v.PostID.Int64,
			User: model.User{
				ID:         v.ID_2.Int64,
				AvatarUrl:  v.AvatarUrl.String,
				Fullname:   v.FullName.String,
				Bio:        v.Bio.String,
				OpenToWork: v.OpenToWork.Bool,
			},
			Content:      v.Content.String,
			ImageUrl:     v.ImageUrl.String,
			LikeCount:    v.LikeCount.Int32,
			ReplyCount:   v.ReplyCount.Int32,
			IsPostAuthor: v.IsPostAuthor.Bool,
			UpdatedAt:    v.UpdatedAt.Time,
		}
	}

	totalPages := int((totalRows + int64(pagination.Limit) - 1) / int64(pagination.Limit))

	paginate := model.PaginationResponse{
		Page:             pagination.Page,
		TotalRows:        totalRows,
		TotalPages:       totalPages,
		CurrentRowsCount: len(data),
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "Success get post comments")
	resp.Data = map[string]any{
		"pagination": paginate,
		"data":       postComments,
	}
	return
}

func (u *PostsUsecase) GetPostCommentReplies(postId, postCommentId int64, pagination model.PaginationRequest) (resp model.Response) {
	offset := (pagination.Page - 1) * pagination.Limit
	data, totalRows, err := u.repository.GetPostCommentReplies(postId, postCommentId, int32(offset), int32(pagination.Limit))

	if err != nil {
		u.log.Errorf("repository.GetPostCommentReplies: %v", err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occurred"),
		}
	}

	postCommentReplies := make([]model.PostCommentReply, len(data))
	for i, v := range data {
		postCommentReplies[i] = model.PostCommentReply{
			ID:            v.ID,
			PostCommentId: v.PostCommentID.Int64,
			User: model.User{
				ID:         v.ID_2.Int64,
				AvatarUrl:  v.AvatarUrl.String,
				Fullname:   v.FullName.String,
				Bio:        v.Bio.String,
				OpenToWork: v.OpenToWork.Bool,
			},
			Content:      v.Content.String,
			ImageUrl:     v.ImageUrl.String,
			LikeCount:    v.LikeCount.Int32,
			IsPostAuthor: v.IsPostAuthor.Bool,
			UpdatedAt:    v.UpdatedAt.Time,
		}
	}

	totalPages := int((totalRows + int64(pagination.Limit) - 1) / int64(pagination.Limit))

	paginate := model.PaginationResponse{
		Page:             pagination.Page,
		TotalRows:        totalRows,
		TotalPages:       totalPages,
		CurrentRowsCount: len(data),
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "Success get post comment replies")
	resp.Data = map[string]any{
		"pagination": paginate,
		"data":       postCommentReplies,
	}
	return
}

func (u *PostsUsecase) UpdatePostLikeCount(postId int64) (resp model.Response) {
	data, err := u.repository.UpdatePostLikeCount(postId)
	if err != nil && err == sql.ErrNoRows {
		return model.Response{
			Status: libs.CustomResponse(http.StatusNotFound, "Data not found"),
		}
	} else if err != nil {
		u.log.Errorf("repository.UpdatePostLikeCount: %v", err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occurred"),
		}
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "Success update post like count")
	resp.Data = map[string]any{
		"id":         data.ID,
		"like_count": data.LikeCount.Int32,
	}
	return
}

func (u *PostsUsecase) ListNewestPostsByUserId(userId int64, pagination model.PaginationRequest) (resp model.Response) {
	offset := (pagination.Page - 1) * pagination.Limit
	data, totalRows, err := u.repository.ListNewestPostsByUserId(userId, int32(offset), int32(pagination.Limit))

	if err != nil {
		u.log.Errorf("repository.ListNewestPostsByUserId (user id %d): %v", userId, err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occurred"),
		}
	}

	totalPages := int((totalRows + int64(pagination.Limit) - 1) / int64(pagination.Limit))

	paginate := model.PaginationResponse{
		Page:             pagination.Page,
		TotalRows:        totalRows,
		TotalPages:       totalPages,
		CurrentRowsCount: len(data),
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "Success get user's posts")
	resp.Data = map[string]any{
		"pagination": paginate,
		"data":       data,
	}
	return
}

func (u *PostsUsecase) ListLikedPostsByUserId(userId int64, pagination model.PaginationRequest) (resp model.Response) {
	offset := (pagination.Page - 1) * pagination.Limit
	data, totalRows, err := u.repository.ListLikedPostsByUserId(userId, int32(offset), int32(pagination.Limit))

	if err != nil {
		u.log.Errorf("repository.ListLikedPostsByUserId (user id %d): %v", userId, err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occurred"),
		}
	}

	totalPages := int((totalRows + int64(pagination.Limit) - 1) / int64(pagination.Limit))

	paginate := model.PaginationResponse{
		Page:             pagination.Page,
		TotalRows:        totalRows,
		TotalPages:       totalPages,
		CurrentRowsCount: len(data),
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "Success get user's posts")
	resp.Data = map[string]any{
		"pagination": paginate,
		"data":       data,
	}
	return
}

func (u *PostsUsecase) ListRepostedPostsByUserId(userId int64, pagination model.PaginationRequest) (resp model.Response) {
	offset := (pagination.Page - 1) * pagination.Limit
	data, totalRows, err := u.repository.ListRepostedPostsByUserId(userId, int32(offset), int32(pagination.Limit))

	if err != nil {
		u.log.Errorf("repository.ListRepostedPostsByUserId (user id %d): %v", userId, err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occurred"),
		}
	}

	totalPages := int((totalRows + int64(pagination.Limit) - 1) / int64(pagination.Limit))

	paginate := model.PaginationResponse{
		Page:             pagination.Page,
		TotalRows:        totalRows,
		TotalPages:       totalPages,
		CurrentRowsCount: len(data),
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "Success get user's posts")
	resp.Data = map[string]any{
		"pagination": paginate,
		"data":       data,
	}
	return
}

func (u *PostsUsecase) InsertPost(imageFile *multipart.FileHeader, props *model.CreatePostRequest) model.Response {
	if imageFile != nil {
		var wg sync.WaitGroup
		objectPath := fmt.Sprintf("users/%d/posts", props.UserId)

		errChan := make(chan error, 1)

		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error

			_, err = u.googleBucket.HandleObjectUpload(imageFile, objectPath)
			if err != nil {
				errChan <- fmt.Errorf("googleBucket.HandleObjectUpload: %v", err)
			}
		}()
		wg.Wait()
		close(errChan)

		if err, ok := <-errChan; ok {
			u.log.Errorf("goroutine error: %v", err)

			return model.Response{
				Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occurred"),
			}
		}
	}

	data, err := u.repository.InsertPost(props)

	if err != nil {
		u.log.Errorf("repository.InsertPost (user id %d): %v", props.UserId, err)

		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occurred"),
		}
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusCreated, "Success create post"),
		Data:   data,
	}
}
