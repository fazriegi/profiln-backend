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

func (u *ProfileUsecase) InsertIssuingOrganization(props *model.CompanyRequest) (resp model.Response) {
	issueOriganization, err := u.repository.InsertIssuingOrganization(props.Name)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertIssuingOrganization %v", err)
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to create issue origanization")
	resp.Data = issueOriganization
	return resp
}

func (u *ProfileUsecase) InsertEmploymentType(props *model.CompanyRequest) (resp model.Response) {
	employmentType, err := u.repository.InsertEmploymentType(props.Name)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertEmploymentType %v", err)
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to create employment type")
	resp.Data = employmentType
	return resp
}

func (u *ProfileUsecase) InsertLocationType(props *model.CompanyRequest) (resp model.Response) {
	locationType, err := u.repository.InsertLocationType(props.Name)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertLocationType %v", err)
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to create location type")
	resp.Data = locationType
	return resp
}

func (u *ProfileUsecase) InsertSchool(props *model.CompanyRequest) (resp model.Response) {
	school, err := u.repository.InsertSchool(props.Name)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertSchool %v", err)
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to create school")
	resp.Data = school
	return resp
}

func (u *ProfileUsecase) InsertSkill(props *model.CompanyRequest) (resp model.Response) {
	skill, err := u.repository.InsertSkill(props.Name)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertSkill %v", err)
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to create skill")
	resp.Data = skill
	return resp
}
