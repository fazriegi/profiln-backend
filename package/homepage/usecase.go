package homepage

import (
	"database/sql"
	"net/http"
	"profiln-be/libs"
	"profiln-be/model"
	repository "profiln-be/package/homepage/repository"
	homepage "profiln-be/package/homepage/repository/sqlc"

	"github.com/sirupsen/logrus"
)

type IHomepageUsecase interface {
	ListPosts(userId int64, pagination model.PaginationRequest) (resp model.Response)
	ListFollowsRecommendation(userId int64, pagination model.PaginationRequest) (resp model.Response)
}

type HomepageUsecase struct {
	repository repository.IHomepageRepository
	log        *logrus.Logger
}

func NewHomepageUsecase(repository repository.IHomepageRepository, log *logrus.Logger) IHomepageUsecase {
	return &HomepageUsecase{
		repository,
		log,
	}
}

func (u *HomepageUsecase) ListPosts(userId int64, pagination model.PaginationRequest) (resp model.Response) {
	offset := (pagination.Page - 1) * pagination.Limit
	posts := []homepage.Post{}
	var (
		totalRows int64
		err       error
	)

	if pagination.OrderBy == "newest" {
		posts, totalRows, err = u.repository.ListPosts(userId, int32(offset), int32(pagination.Limit))

		if err != nil {
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

			u.log.Errorf("repository.ListPosts: %v", err)
			return
		}

	} else if pagination.OrderBy == "following" {
		posts, totalRows, err = u.repository.ListPostsByFollowing(userId, int32(offset), int32(pagination.Limit))

		if err != nil {
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

			u.log.Errorf("repository.ListPostsByFollowing: %v", err)
			return
		}
	} else if pagination.OrderBy == "popular" {
		posts, totalRows, err = u.repository.ListPopularPosts(userId, int32(offset), int32(pagination.Limit))

		if err != nil {
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

			u.log.Errorf("repository.ListPopularPosts: %v", err)
			return
		}
	}

	user, err := u.repository.GetUserById(userId)
	if err != nil && err != sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

		u.log.Errorf("repository.GetUserById: %v", err)
		return
	}

	userData := model.User{
		ID:         user.ID,
		AvatarUrl:  user.AvatarUrl.String,
		Fullname:   user.FullName,
		Bio:        user.Bio.String,
		OpenToWork: user.OpenToWork.Bool,
	}

	data := make([]model.ListPostsResponse, len(posts))
	for i, v := range posts {
		data[i] = model.ListPostsResponse{
			ID:           v.ID,
			User:         userData,
			Content:      v.Content.String,
			ImageUrl:     v.ImageUrl.String,
			LikeCount:    v.LikeCount.Int32,
			CommentCount: v.CommentCount.Int32,
			RepostCount:  v.RepostCount.Int32,
			UpdatedAt:    v.UpdatedAt.Time,
		}
	}

	totalPages := int(totalRows / int64(pagination.Limit))
	if totalRows%int64(pagination.Limit) > 0 {
		totalPages++
	}

	paginate := model.PaginationResponse{
		Page:             pagination.Page,
		TotalRows:        totalRows,
		TotalPages:       totalPages,
		CurrentRowsCount: len(data),
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "Success fetch posts")
	resp.Data = map[string]any{
		"pagination": paginate,
		"data":       data,
	}

	return
}

func (u *HomepageUsecase) ListFollowsRecommendation(userId int64, pagination model.PaginationRequest) (resp model.Response) {
	offset := (pagination.Page - 1) * pagination.Limit
	users, totalRows, err :=
		u.repository.GetFollowsRecommendationForUserId(userId, int32(offset), int32(pagination.Limit))

	if err != nil {
		u.log.Errorf("repository.GetFollowsRecommendationForUserId: %v", err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occurred"),
		}
	}

	data := make([]model.User, len(users))
	for i, v := range users {
		data[i] = model.User{
			ID:         v.ID,
			AvatarUrl:  v.AvatarUrl.String,
			Fullname:   v.FullName,
			Bio:        v.Bio.String,
			OpenToWork: v.OpenToWork.Bool,
		}
	}

	totalPages := int((totalRows + int64(pagination.Limit) - 1) / int64(pagination.Limit))

	paginate := model.PaginationResponse{
		Page:             pagination.Page,
		TotalRows:        totalRows,
		TotalPages:       totalPages,
		CurrentRowsCount: len(data),
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "Success get follows recommendations")
	resp.Data = map[string]any{
		"pagination": paginate,
		"data":       data,
	}
	return
}
