package profile

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"net/http"

	"profiln-be/libs"
	"profiln-be/model"
	repository "profiln-be/package/profile/repository"

	profileSqlc "profiln-be/package/profile/repository/sqlc"

	"github.com/sirupsen/logrus"
)

type IProfileUsecase interface {
	InsertCompany(props *model.CompanyRequest) (resp model.Response)
	InsertIssuingOrganization(props *model.IssuingOrganizationRequest) (resp model.Response)
	InsertUserDetail(props *model.UserDetailRequest, id int64) (resp model.Response)
	InsertUserDetailAbout(props *model.UserDetailAboutRequest, id int64) (resp model.Response)
	InsertEducation(props *model.EducationRequest, id int64) (resp model.Response)
	InsertWorkExperience(props *model.WorkExperienceRequest, id int64) (resp model.Response)
	InsertCertificate(props *model.CertificateRequest, id int64) (resp model.Response)
	InsertUserSkill(props *model.UserSkillRequest, id int64) (resp model.Response)
	GetSkills(pagination model.PaginationRequest) (resp model.Response)
	UpdateProfile(imageFile *multipart.FileHeader, props *model.UpdateProfileRequest) (resp model.Response)
	UpdateAboutMe(userId int64, aboutMe string) (resp model.Response)
	UpdateUserCertificate(userId int64, props *model.UpdateCertificate) (resp model.Response)
	UpdateUserInformation(props *model.UpdateUserInformation) (resp model.Response)
	UpdateUserEducation(imageFile *multipart.FileHeader, props *model.UpdateEducationRequest) (resp model.Response)
}

type ProfileUsecase struct {
	repository   repository.IProfileRepository
	log          *logrus.Logger
	googleBucket libs.IGoogleBucket
}

func NewProfileUsecase(repository repository.IProfileRepository, log *logrus.Logger, googleBucket libs.IGoogleBucket) IProfileUsecase {
	return &ProfileUsecase{
		repository,
		log,
		googleBucket,
	}
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
		// FinishDate:   sql.NullTime{Time: time.Time{}, Valid: false},
		FinishDate: sql.NullTime{Time: props.FinishDate.Time, Valid: true},
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
	insertAboutParams := profileSqlc.InsertUserDetailAboutParams{
		About:  sql.NullString{String: props.About, Valid: true},
		UserID: sql.NullInt64{Int64: id, Valid: true},
	}

	_, err := u.repository.GetUserById(id)
	if err != nil && err == sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusNotFound, "User not found")

		return
	} else if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

		u.log.Errorf("repository.GetUserById: %v", err)
		return
	}

	about, err := u.repository.InsertUserDetailAbout(insertAboutParams)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("repository.InsertUserDetailAbout %v", err)
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to create user about")
	resp.Data = about.About.String
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

func (u *ProfileUsecase) GetSkills(pagination model.PaginationRequest) (resp model.Response) {
	offset := (pagination.Page - 1) * pagination.Limit

	skills, totalRows, err := u.repository.GetSkills(int32(offset), int32(pagination.Limit))
	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")
		u.log.Errorf("repository.GetSkills: %v", err)
		return
	}

	data := make([]model.GetSkillsResponse, len(skills))
	for i, v := range skills {
		data[i] = model.GetSkillsResponse{
			ID:   v.ID,
			Name: v.Name,
		}
	}

	totalPages := int((totalRows + int64(pagination.Limit) - 1) / int64(pagination.Limit))

	paginate := model.PaginationResponse{
		Page:             pagination.Page,
		TotalRows:        totalRows,
		TotalPages:       totalPages,
		CurrentRowsCount: len(data),
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success fetch skills"),
		Data: map[string]any{
			"pagination": paginate,
			"data":       data,
		},
	}
}

func (u *ProfileUsecase) UpdateProfile(imageFile *multipart.FileHeader, props *model.UpdateProfileRequest) (resp model.Response) {
	var (
		avatarUrl string
		err       error
	)

	avatarUrl, err = u.repository.GetUserAvatarById(props.UserId)
	if err != nil && err != sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")
		u.log.Errorf("repository.GetUserAvatarById: %v", err)
		return
	}

	if imageFile != nil {
		objectPath := fmt.Sprintf("users/%d/avatar", props.UserId)

		avatarUrl, err = u.googleBucket.HandleObjectUpload(imageFile, avatarUrl, objectPath)
		if err != nil {
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")
			u.log.Errorf("googleBucket.HandleObjectUpload: %v", err)
			return
		}
	}

	err = u.repository.UpdateProfile(avatarUrl, props)
	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")
		u.log.Errorf("repository.UpdateProfile: %v", err)
		return
	}

	responseData := model.UpdateProfileResponse{
		UserId:          props.UserId,
		Fullname:        props.Fullname,
		AvatarUrl:       avatarUrl,
		HidePhoneNumber: props.HidePhoneNumber,
		MainSkills:      props.MainSkills,
		PhoneNumber:     props.PhoneNumber,
		Gender:          props.Gender,
		SocialLinks:     props.SocialLinks,
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success edit profile"),
		Data:   responseData,
	}
}

func (u *ProfileUsecase) UpdateAboutMe(userId int64, aboutMe string) (resp model.Response) {
	err := u.repository.UpdateAboutMe(userId, aboutMe)
	if err != nil && err == sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusNotFound, "Data not found")
		return
	} else if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")
		u.log.Errorf("repository.UpdateAboutMe: %v", err)
		return
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success update user's about"),
		Data: map[string]any{
			"user_id": userId,
			"about":   aboutMe,
		},
	}
}

func (u *ProfileUsecase) UpdateUserCertificate(userId int64, props *model.UpdateCertificate) (resp model.Response) {
	err := u.repository.UpdateUserCertificate(userId, props)
	if err != nil && err == sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusNotFound, "Data not found")
		return
	} else if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")
		u.log.Errorf("repository.UpdateAboutMe: %v", err)
		return
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success update user's certificate"),
		Data:   props,
	}
}

func (u *ProfileUsecase) UpdateUserInformation(props *model.UpdateUserInformation) (resp model.Response) {
	err := u.repository.UpdateUserInformation(props)
	if err != nil && err == sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusNotFound, "Data not found")
		return
	} else if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")
		u.log.Errorf("repository.UpdateAboutMe: %v", err)
		return
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success update user's information"),
		Data:   props,
	}
}
