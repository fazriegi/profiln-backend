package profile

import (
	"database/sql"
	"net/http"

	// "os"
	// "time"

	"profiln-be/libs"
	email "profiln-be/libs/email"
	"profiln-be/model"
	repository "profiln-be/package/profile/repository"

	profileSqlc "profiln-be/package/profile/repository/sqlc"

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

func (u *ProfileUsecase) InsertIssuingOrganization(props *model.IssuingOrganizationRequest) (resp model.Response) {
	issueOriganization, err := u.repository.InsertIssuingOrganization(props.Name)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertIssuingOrganization %v", err)
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to create issue origanization")
	resp.Data = issueOriganization
	return resp
}

func (u *ProfileUsecase) InsertEmploymentType(props *model.EmploymentTypeRequest) (resp model.Response) {
	employmentType, err := u.repository.InsertEmploymentType(props.Name)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertEmploymentType %v", err)
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to create employment type")
	resp.Data = employmentType
	return resp
}

func (u *ProfileUsecase) InsertLocationType(props *model.LocationTypeRequest) (resp model.Response) {
	locationType, err := u.repository.InsertLocationType(props.Name)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertLocationType %v", err)
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to create location type")
	resp.Data = locationType
	return resp
}

func (u *ProfileUsecase) InsertSchool(props *model.SchoolRequest) (resp model.Response) {
	school, err := u.repository.InsertSchool(props.Name)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertSchool %v", err)
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to create school")
	resp.Data = school
	return resp
}

func (u *ProfileUsecase) InsertSkill(props *model.SkillRequest) (resp model.Response) {
	skill, err := u.repository.InsertSkill(props.Name)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertSkill %v", err)
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to create skill")
	resp.Data = skill
	return resp
}

func (u *ProfileUsecase) InsertUserSkill(props *model.UserSkillRequest, id int64) (resp model.Response) {
	skill, err := u.repository.InsertSkill(props.Skills)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertSkill %v", err)
	}

	userSkillParams := profileSqlc.InsertUserSkillParams{
		UserID:    sql.NullInt64{Int64: id, Valid: true},
		SkillID:   sql.NullInt64{Int64: skill.ID, Valid: true},
		MainSkill: sql.NullBool{Bool: false, Valid: true},
	}

	userSkill, err := u.repository.InsertUserSkill(userSkillParams)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertUserSkill %v", err)
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to create skill")
	resp.Data = userSkill
	return resp
}

func (u *ProfileUsecase) InsertCertificate(props *model.CertificateRequest, id int64) (resp model.Response) {
	certificateParams := profileSqlc.InsertCertificateParams{
		UserID:                sql.NullInt64{Int64: id, Valid: true},
		Name:                  sql.NullString{String: props.Name, Valid: true},
		IssuingOrganizationID: sql.NullInt64{Int64: props.IssuingOrganizationID, Valid: true},
		IssueDate:             sql.NullTime{Time: props.IssueDate.Time, Valid: true},
		ExpirationDate:        sql.NullTime{Time: props.ExpirationDate.Time, Valid: true},
		CredentialID:          sql.NullString{String: props.CredentialID, Valid: true},
		Url:                   sql.NullString{String: props.Url, Valid: true},
	}

	certificate, err := u.repository.InsertCertificate(certificateParams)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertCertificate %v", err)
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to create certificate")
	resp.Data = certificate
	return resp
}

func (u *ProfileUsecase) InsertWorkExperience(props *model.WorkExperienceRequest, id int64) (resp model.Response) {
	workExperienceParams := profileSqlc.InsertWorkExperienceParams{
		UserID:           sql.NullInt64{Int64: id, Valid: true},
		JobTitle:         sql.NullString{String: props.JobTitle, Valid: true},
		CompanyID:        sql.NullInt64{Int64: props.CompanyID, Valid: true},
		EmploymentTypeID: sql.NullInt16{Int16: props.EmploymentTypeID, Valid: true},
		Location:         sql.NullString{String: props.Location, Valid: true},
		LocationTypeID:   sql.NullInt16{Int16: props.LocationTypeID, Valid: true},
		StartDate:        sql.NullTime{Time: props.StartDate.Time, Valid: true},
		FinishDate:       sql.NullTime{Time: props.FinishDate.Time, Valid: true},
		Description:      sql.NullString{String: props.Description, Valid: true},
	}

	workExperience, err := u.repository.InsertWorkExperience(workExperienceParams)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertWorkExperience %v", err)
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to create work experience")
	resp.Data = workExperience
	return resp
}

func (u *ProfileUsecase) InsertEducation(props *model.EducationRequest, id int64) (resp model.Response) {
	educationParams := profileSqlc.InsertEducationParams{
		UserID:       sql.NullInt64{Int64: id, Valid: true},
		SchoolID:     sql.NullInt64{Int64: props.SchoolID, Valid: true},
		Degree:       sql.NullString{String: props.Degree, Valid: true},
		FieldOfStudy: sql.NullString{String: props.FieldOfStudy, Valid: true},
		Gpa:          sql.NullString{String: props.Gpa, Valid: true},
		StartDate:    sql.NullTime{Time: props.StartDate.Time, Valid: true},
		FinishDate:   sql.NullTime{Time: props.FinishDate.Time, Valid: true},
	}

	education, err := u.repository.InsertEducation(educationParams)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertEducation %v", err)
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to create education")
	resp.Data = education
	return resp
}

func (u *ProfileUsecase) InsertUserDetailAbout(props *model.UserDetailAboutRequest, id int64) (resp model.Response) {
	aboutParams := profileSqlc.InsertUserDetailAboutParams{
		About:  sql.NullString{String: props.About, Valid: true},
		UserID: sql.NullInt64{Int64: id, Valid: true},
	}

	err := u.repository.InsertUserDetailAbout(aboutParams)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertEducation %v", err)
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to create user about")
	return resp
}

func (u *ProfileUsecase) InsertUserDetail(props *model.UserDetailRequest, id int64) (resp model.Response) {
	userDetailParams := profileSqlc.InsertUserDetailParams{
		UserID:      sql.NullInt64{Int64: id, Valid: true},
		PhoneNumber: sql.NullString{String: props.PhoneNumber, Valid: true},
		Gender:      sql.NullString{String: props.Gender, Valid: true},
	}

	userDetail, err := u.repository.InsertUserDetail(userDetailParams)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertUserDetail %v", err)
	}

	avatarParams := profileSqlc.InsertUserAvatarParams{
		AvatarUrl: sql.NullString{String: "", Valid: true},
		ID:        id,
	}

	err = u.repository.InsertUserAvatar(avatarParams)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertUserAvatar %v", err)
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to create user about")
	resp.Data = userDetail
	return resp
}
