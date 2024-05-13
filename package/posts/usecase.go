package posts

import (
	"net/http"
	"profiln-be/libs"
	"profiln-be/model"
	repository "profiln-be/package/posts/repository"

	"github.com/sirupsen/logrus"
)

type IPostsUsecase interface {
	InsertReportedPost(userId int64, data *model.ReportPostRequest) (resp model.Response)
}

type PostsUsecase struct {
	repository repository.IPostsRepository
	log        *logrus.Logger
}

func NewPostsUsecase(repository repository.IPostsRepository, log *logrus.Logger) IPostsUsecase {
	return &PostsUsecase{
		repository,
		log,
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
