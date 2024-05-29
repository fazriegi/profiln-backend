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
	InsertPost(imageFile []*multipart.FileHeader, props *model.CreatePostRequest) model.Response
	UpdatePost(imageFiles []*multipart.FileHeader, props *model.UpdatePostRequest) model.Response
	DeletePost(userId, postId int64) model.Response
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

	resp.Status = libs.CustomResponse(http.StatusOK, "Success get detail post")
	resp.Data = data
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

func (u *PostsUsecase) InsertPost(imageFiles []*multipart.FileHeader, props *model.CreatePostRequest) model.Response {
	if imageFiles != nil {
		var wg sync.WaitGroup
		objectPath := fmt.Sprintf("users/%d/posts", props.UserId)

		errChan := make(chan error, len(imageFiles))
		urlChan := make(chan string, len(imageFiles))

		// Loop through the imageFiles
		for _, file := range imageFiles {
			wg.Add(1)
			file := file

			// Handle object uploads to gcloud storage for each file asynchronously
			go func(file *multipart.FileHeader) {
				defer wg.Done()
				objectUrl, err := u.googleBucket.HandleObjectUpload(file, objectPath)

				if err != nil {
					errChan <- fmt.Errorf("googleBucket.HandleObjectUpload (user id: %d): %v", props.UserId, err)
					return
				}

				urlChan <- objectUrl

			}(file)
		}

		wg.Wait()
		close(errChan)
		close(urlChan)

		// Loop through error channel and check if any error occurred
		for err := range errChan {
			if err != nil {
				u.log.Error(err)

				return model.Response{
					Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
				}
			}
		}

		for url := range urlChan {
			props.ImageUrls = append(props.ImageUrls, url)
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

func (u *PostsUsecase) UpdatePost(imageFiles []*multipart.FileHeader, props *model.UpdatePostRequest) model.Response {
	var (
		err error
	)

	// Check if user post exists
	currentPost, err := u.repository.GetPostById(props.ID)
	if err == sql.ErrNoRows {
		return model.Response{
			Status: libs.CustomResponse(http.StatusNotFound, "Data not found"),
		}
	}

	if currentPost.User.ID != props.UserId {
		return model.Response{
			Status: libs.CustomResponse(http.StatusUnauthorized, "Unauthorized"),
		}
	}

	// Get all current post file urls
	currentObjectUrls, err := u.repository.GetPostImagesUrl(props.ID)
	if err != nil {
		u.log.Errorf("repository.GetPostImagesUrl (user id: %d): %v", props.UserId, err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	props.ImageUrls = currentObjectUrls

	if imageFiles != nil {
		var wg sync.WaitGroup
		objectPath := fmt.Sprintf("users/%d/posts", props.UserId)

		errChan := make(chan error, len(imageFiles))
		urlChan := make(chan string, len(imageFiles))

		// Loop through the imageFiles
		for _, file := range imageFiles {
			wg.Add(1)
			file := file

			// Handle object uploads to gcloud storage for each file asynchronously
			go func(file *multipart.FileHeader) {
				defer wg.Done()
				objectUrl, err := u.googleBucket.HandleObjectUpload(file, objectPath)

				if err != nil {
					errChan <- fmt.Errorf("googleBucket.HandleObjectUpload (user id: %d): %v", props.UserId, err)
					return
				}

				urlChan <- objectUrl

			}(file)
		}

		wg.Wait()
		close(errChan)
		close(urlChan)

		// Loop through error channel and check if any error occurred
		for err := range errChan {
			if err != nil {
				u.log.Error(err)

				return model.Response{
					Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
				}
			}
		}

		props.ImageUrls = []string{}
		for url := range urlChan {
			props.ImageUrls = append(props.ImageUrls, url)
		}
	}

	err = u.repository.UpdatePostById(props)
	if err != nil {
		// Delete uploaded objects
		errObjectDelete := u.googleBucket.HandleObjectDeletion(props.ImageUrls...)
		if errObjectDelete != nil {
			u.log.Errorf("googleBucket.HandleObjectDeletion (user id: %d): %v", props.UserId, errObjectDelete)
		}

		if err == sql.ErrNoRows {
			return model.Response{
				Status: libs.CustomResponse(http.StatusNotFound, "Data not found"),
			}
		}

		u.log.Errorf("repository.UpdatePostById (user id: %d): %v", props.UserId, err)

		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	// If previous objects exists, delete it from gcloud storage
	if len(currentObjectUrls) > 1 && imageFiles != nil {
		err := u.googleBucket.HandleObjectDeletion(currentObjectUrls...)
		if err != nil {
			u.log.Errorf("googleBucket.HandleObjectDeletion (user id: %d): %v", props.UserId, err)
		}
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success update post"),
		Data:   props,
	}
}

func (u *PostsUsecase) DeletePost(userId, postId int64) model.Response {
	// Check if user post exists
	currentPost, err := u.repository.GetPostById(postId)
	if err == sql.ErrNoRows {
		return model.Response{
			Status: libs.CustomResponse(http.StatusNotFound, "Data not found"),
		}
	}

	if currentPost.User.ID != userId {
		return model.Response{
			Status: libs.CustomResponse(http.StatusUnauthorized, "Unauthorized"),
		}
	}

	// Get all current post file urls
	currentObjectUrls, err := u.repository.GetPostImagesUrl(postId)
	if err != nil {
		u.log.Errorf("repository.GetPostImagesUrl (post id: %d): %v", postId, err)

		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	err = u.repository.DeletePost(postId)
	if err != nil {
		u.log.Errorf("repository.DeletePost(%d): %v", postId, err)

		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	if len(currentObjectUrls) > 1 {
		err := u.googleBucket.HandleObjectDeletion(currentObjectUrls...)
		if err != nil {
			u.log.Errorf("googleBucket.HandleObjectDeletion (user id: %d): %v", userId, err)
		}
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success delete post"),
	}
}
