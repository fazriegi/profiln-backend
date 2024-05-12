package profile

import (
	// "database/sql"
	// "fmt"
	"net/http"
	// "os"
	// "time"

	"profiln-be/libs"
	email "profiln-be/libs/email"
	"profiln-be/model"
	repository "profiln-be/package/profile/repository"

	// profileSqlc "profiln-be/package/profile/repository/sqlc"

	"github.com/sirupsen/logrus"
)

type ProfileUsecase struct {
	repository repository.IProfileRepository
	email      email.IEmail
	log        *logrus.Logger
}

func (u *ProfileUsecase) InsertCompany(props *model.CompanyRequest) (resp model.Response) {
	company, err := u.repository.InsertCompany(props.Name)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertCompany %v", err)
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to create company")
	resp.Data = company
	return resp
}
