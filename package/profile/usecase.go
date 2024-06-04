package profile

import (
	"database/sql"
	"fmt"
	"net/http"

	"profiln-be/libs"
	"profiln-be/model"
	repository "profiln-be/package/profile/repository"

	"github.com/sirupsen/logrus"
)

type IProfileUsecase interface {
	InsertUserProfile(imageFileNames []string, props *model.AddProfileRequest) model.Response
	InsertUserWorkExperience(fileNames []string, props *model.WorkExperience) model.Response
	InsertUserEducation(filenames []string, props *model.Education) model.Response
	InsertUserCertificate(props *model.Certificate) model.Response
	InsertUserSkills(props *model.AddUserSkill) model.Response
	UpdateProfile(imageFileNames []string, props *model.UpdateProfileRequest) (resp model.Response)
	UpdateAboutMe(userId int64, aboutMe string) (resp model.Response)
	UpdateUserCertificate(userId int64, props *model.Certificate) (resp model.Response)
	UpdateUserInformation(props *model.UpdateUserInformation) (resp model.Response)
	UpdateUserEducation(fileNames []string, props *model.Education) (resp model.Response)
	UpdateUserWorkExperience(fileNames []string, props *model.WorkExperience) (resp model.Response)
	AddUserOpenToWork(props *model.OpenToWork) model.Response
	GetUserProfile(userId int64) model.Response
	GetWorkExperiencesByUserId(userId int64, pagination model.PaginationRequest) model.Response
	GetEducationsByUserId(userId int64, pagination model.PaginationRequest) model.Response
	GetCertificatesByUserId(userId int64, pagination model.PaginationRequest) model.Response
	GetFollowedUsersByUserId(userId int64, pagination model.PaginationRequest) model.Response
	GetUserBasicInformation(userId int64) model.Response
	DeleteUserOpenToWork(userId int64) model.Response
	DeleteUserWorkExperienceById(userId, workExperienceId int64) model.Response
	DeleteUserEducationById(userId, educationId int64) model.Response
	DeleteUserCertificateById(userId, educationId int64) model.Response
	FollowUser(userId, targetUserId int64) model.Response
	UnfollowUser(userId, targetUserId int64) model.Response
}

type ProfileUsecase struct {
	repository   repository.IProfileRepository
	log          *logrus.Logger
	googleBucket libs.IGoogleBucket
	fs           libs.IFileSystem
}

func NewProfileUsecase(repository repository.IProfileRepository, log *logrus.Logger, googleBucket libs.IGoogleBucket, fs libs.IFileSystem) IProfileUsecase {
	return &ProfileUsecase{
		repository,
		log,
		googleBucket,
		fs,
	}
}

func (u *ProfileUsecase) InsertUserSkills(props *model.AddUserSkill) (resp model.Response) {
	err := u.repository.BatchInsertUserSkills(props.UserId, props.Skills)
	if err != nil {
		u.log.Errorf("repository.BatchInsertUserSkills (user id: %d): %v", props.UserId, err)

		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusCreated, "Success add user's skills"),
		Data:   props,
	}
}

func (u *ProfileUsecase) InsertUserCertificate(props *model.Certificate) model.Response {
	data, err := u.repository.InsertUserCertificate(props)
	if err != nil {
		u.log.Errorf("repository.InsertUserCertificate (user id: %d): %v", props.UserId, err)

		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusCreated, "Success add user's certificate"),
		Data:   data,
	}
}

func (u *ProfileUsecase) UpdateProfile(imageFileNames []string, props *model.UpdateProfileRequest) (resp model.Response) {
	var (
		avatarUrl string
		err       error
	)

	currentAvatarUrl, err := u.repository.GetUserAvatarById(props.UserId)
	if err != nil && err != sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")
		u.log.Errorf("repository.GetUserAvatarById: %v", err)
		return
	}

	if len(imageFileNames) > 0 {
		objectPath := fmt.Sprintf("users/%d/avatar", props.UserId)

		defer func() {
			filePath := fmt.Sprintf("./storage/temp/users/%d/files/%s", props.UserId, imageFileNames[0])

			if err := u.fs.RemoveFile(filePath); err != nil {
				u.log.Errorf("fileSystem.RemoveFile: %v", err)
			}
		}()

		urls, err := u.googleBucket.HandleObjectUploads(props.UserId, objectPath, imageFileNames[0])
		if err != nil {
			u.log.Errorf("googleBucket.HandleObjectUploads: %v", err)

			return model.Response{
				Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
			}
		}

		avatarUrl = urls[0]
	}

	err = u.repository.UpdateProfile(avatarUrl, props)
	if err != nil {
		errObjectDelete := u.googleBucket.HandleObjectDeletion(avatarUrl)
		if errObjectDelete != nil {
			u.log.Errorf("googleBucket.HandleObjectDeletion: %v", errObjectDelete)
		}

		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")
		u.log.Errorf("repository.UpdateProfile: %v", err)
		return
	}

	if currentAvatarUrl != "" {
		err := u.googleBucket.HandleObjectDeletion(currentAvatarUrl)
		if err != nil {
			u.log.Errorf("googleBucket.HandleObjectDeletion: %v", err)
		}
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

func (u *ProfileUsecase) UpdateUserCertificate(userId int64, props *model.Certificate) (resp model.Response) {
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

func (u *ProfileUsecase) UpdateUserEducation(fileNames []string, props *model.Education) (resp model.Response) {
	var (
		err error
	)

	// Check if user education exists
	_, err = u.repository.GetEducationById(props.ID)
	if err == sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusNotFound, "Data not found")
		return
	}

	// Get all current education file urls
	currentObjectUrls, err := u.repository.GetUserEducationFileURLs(props.ID)
	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")
		u.log.Errorf("repository.GetUserEducationFileURLs (user id: %d): %v", props.UserId, err)
		return
	}

	props.FileURLs = currentObjectUrls

	if len(fileNames) > 0 {
		defer func() {
			for _, fileName := range fileNames {
				filePath := fmt.Sprintf("./storage/temp/users/%d/files/%s", props.UserId, fileName)

				if err := u.fs.RemoveFile(filePath); err != nil {
					u.log.Errorf("fileSystem.RemoveFile: %v", err)
				}
			}
		}()

		objectPath := fmt.Sprintf("users/%d/educations/files", props.UserId)

		urls, err := u.googleBucket.HandleObjectUploads(props.UserId, objectPath, fileNames...)
		if err != nil {
			u.log.Errorf("googleBucket.HandleObjectUploads: %v", err)

			return model.Response{
				Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
			}
		}

		props.FileURLs = urls
	}

	err = u.repository.UpdateUserEducation(props)
	if err != nil {
		// Delete uploaded objects
		errObjectDelete := u.googleBucket.HandleObjectDeletion(props.FileURLs...)
		if errObjectDelete != nil {
			u.log.Errorf("googleBucket.HandleObjectDeletion (user id: %d): %v", props.UserId, errObjectDelete)
		}

		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")
		u.log.Errorf("repository.UpdateUserEducation (user id: %d): %v", props.UserId, err)
		return
	}

	// If previous objects exists, delete it from gcloud storage
	if len(currentObjectUrls) > 0 {
		err := u.googleBucket.HandleObjectDeletion(currentObjectUrls...)
		if err != nil {
			u.log.Errorf("googleBucket.HandleObjectDeletionc (user id: %d): %v", props.UserId, err)
		}
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success update user's education"),
		Data:   props,
	}
}

func (u *ProfileUsecase) UpdateUserWorkExperience(fileNames []string, props *model.WorkExperience) (resp model.Response) {
	var (
		err error
	)

	// Check if user work experience exists
	_, err = u.repository.GetWorkExperienceById(props.ID)
	if err == sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusNotFound, "Data not found")
		return
	}

	// Get all current work experience file urls
	currentObjectUrls, err := u.repository.GetWorkExperienceFileURLs(props.ID)
	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")
		u.log.Errorf("repository.GetWorkExperienceFileURLs (user id: %d): %v", props.UserId, err)
		return
	}

	props.FileURLs = currentObjectUrls

	if fileNames != nil {
		objectPath := fmt.Sprintf("users/%d/work-experiences/files", props.UserId)

		defer func() {
			for _, fileName := range fileNames {
				filePath := fmt.Sprintf("./storage/temp/users/%d/files/%s", props.UserId, fileName)

				if err := u.fs.RemoveFile(filePath); err != nil {
					u.log.Errorf("fileSystem.RemoveFile: %v", err)
				}
			}
		}()

		urls, err := u.googleBucket.HandleObjectUploads(props.UserId, objectPath, fileNames...)
		if err != nil {
			u.log.Errorf("googleBucket.HandleObjectUploads: %v", err)

			return model.Response{
				Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
			}
		}

		props.FileURLs = urls
	}

	err = u.repository.UpdateUserWorkExperience(props)
	if err != nil {
		// Delete uploaded objects
		errObjectDelete := u.googleBucket.HandleObjectDeletion(props.FileURLs...)
		if errObjectDelete != nil {
			u.log.Errorf("googleBucket.HandleObjectDeletion (user id: %d): %v", props.UserId, errObjectDelete)
		}

		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")
		u.log.Errorf("repository.UpdateUserWorkExperience (user id: %d): %v", props.UserId, err)
		return
	}

	// If previous objects exists, delete it from gcloud storage
	if len(currentObjectUrls) > 0 {
		err := u.googleBucket.HandleObjectDeletion(currentObjectUrls...)
		if err != nil {
			u.log.Errorf("googleBucket.HandleObjectDeletion (user id: %d): %v", props.UserId, err)
		}
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success update user's work experience"),
		Data:   props,
	}
}

func (u *ProfileUsecase) GetUserProfile(userId int64) model.Response {
	data, err := u.repository.GetUserProfile(userId)
	if err != nil && err == sql.ErrNoRows {
		return model.Response{
			Status: libs.CustomResponse(http.StatusNotFound, "Data not found"),
		}
	} else if err != nil {
		u.log.Errorf("repository.GetUserProfile(%d): %v", userId, err)

		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success fetch user profile"),
		Data:   data,
	}
}

func (u *ProfileUsecase) GetWorkExperiencesByUserId(userId int64, pagination model.PaginationRequest) model.Response {
	offset := (pagination.Page - 1) * pagination.Limit

	data, totalRows, err := u.repository.GetWorkExperiencesByUserId(userId, int32(offset), int32(pagination.Limit))
	if err != nil {
		u.log.Errorf("repository.GetWorkExperiencesByUserId(%d): %v", userId, err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
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
		Status: libs.CustomResponse(http.StatusOK, "Success fetch user work experiences"),
		Data: map[string]any{
			"pagination": paginate,
			"data":       data,
		},
	}
}

func (u *ProfileUsecase) GetEducationsByUserId(userId int64, pagination model.PaginationRequest) model.Response {
	offset := (pagination.Page - 1) * pagination.Limit

	data, totalRows, err := u.repository.GetEducationsByUserId(userId, int32(offset), int32(pagination.Limit))
	if err != nil {
		u.log.Errorf("repository.GetEducationsByUserId(%d): %v", userId, err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
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
		Status: libs.CustomResponse(http.StatusOK, "Success fetch user educations"),
		Data: map[string]any{
			"pagination": paginate,
			"data":       data,
		},
	}
}

func (u *ProfileUsecase) GetCertificatesByUserId(userId int64, pagination model.PaginationRequest) model.Response {
	offset := (pagination.Page - 1) * pagination.Limit

	data, totalRows, err := u.repository.GetCertificatesByUserId(userId, int32(offset), int32(pagination.Limit))
	if err != nil {
		u.log.Errorf("repository.GetCertificatesByUserId(%d): %v", userId, err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
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
		Status: libs.CustomResponse(http.StatusOK, "Success fetch user certificates"),
		Data: map[string]any{
			"pagination": paginate,
			"data":       data,
		},
	}
}

func (u *ProfileUsecase) GetFollowedUsersByUserId(userId int64, pagination model.PaginationRequest) model.Response {
	offset := (pagination.Page - 1) * pagination.Limit

	data, totalRows, err := u.repository.GetFollowedUsersByUserId(userId, int32(offset), int32(pagination.Limit))
	if err != nil {
		u.log.Errorf("repository.GetFollowedUsersByUserId(%d): %v", userId, err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
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
		Status: libs.CustomResponse(http.StatusOK, "Success fetch followed users"),
		Data: map[string]any{
			"pagination": paginate,
			"data":       data,
		},
	}
}

func (u *ProfileUsecase) GetUserBasicInformation(userId int64) model.Response {
	data, err := u.repository.GetUserById(userId)
	if err != nil && err == sql.ErrNoRows {
		return model.Response{
			Status: libs.CustomResponse(http.StatusNotFound, "Data not found"),
		}
	} else if err != nil {
		u.log.Errorf("repository.GetUserById(%d): %v", userId, err)

		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success fetch user basic information"),
		Data:   data,
	}
}

func (u *ProfileUsecase) AddUserOpenToWork(props *model.OpenToWork) model.Response {
	props.OpenToWork = true
	err := u.repository.AddUserOpenToWork(props)
	if err != nil {
		u.log.Errorf("repository.AddUserOpenToWork: %v", err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	return model.Response{
		Data:   props,
		Status: libs.CustomResponse(http.StatusOK, "Success update user open to work"),
	}
}

func (u *ProfileUsecase) DeleteUserOpenToWork(userId int64) model.Response {
	err := u.repository.DeleteUserOpenToWork(userId)
	if err != nil {
		u.log.Errorf("repository.DeleteUserOpenToWork: %v", err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success delete user open to work"),
	}
}

func (u *ProfileUsecase) DeleteUserWorkExperienceById(userId, workExperienceId int64) model.Response {
	fileUrls, err := u.repository.GetWorkExperienceFileURLs(workExperienceId)
	if err != nil {
		u.log.Errorf("repository.GetWorkExperienceFileURLs: %v", err)

		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	err = u.repository.DeleteUserWorkExperienceById(userId, workExperienceId)
	if err != nil {
		u.log.Errorf("repository.DeleteUserWorkExperienceById: %v", err)

		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	if len(fileUrls) > 0 {
		err := u.googleBucket.HandleObjectDeletion(fileUrls...)
		if err != nil {
			u.log.Errorf("googleBucket.HandleObjectDeletion: %v", err)
		}
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success delete user work experience"),
	}
}

func (u *ProfileUsecase) DeleteUserEducationById(userId, educationId int64) model.Response {
	fileUrls, err := u.repository.GetUserEducationFileURLs(educationId)
	if err != nil {
		u.log.Errorf("repository.GetUserEducationFileURLs: %v", err)

		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	err = u.repository.DeleteUserEducationById(userId, educationId)
	if err != nil {
		u.log.Errorf("repository.DeleteUserEducationById: %v", err)

		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	if len(fileUrls) > 0 {
		err := u.googleBucket.HandleObjectDeletion(fileUrls...)
		if err != nil {
			u.log.Errorf("googleBucket.HandleObjectDeletion: %v", err)
		}
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success delete user education"),
	}
}

func (u *ProfileUsecase) DeleteUserCertificateById(userId, certificateId int64) model.Response {
	err := u.repository.DeleteUserCertificateById(userId, certificateId)
	if err != nil {
		u.log.Errorf("repository.DeleteUserCertificateById: %v", err)

		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success delete user certificate"),
	}
}

func (u *ProfileUsecase) FollowUser(userId, targetUserId int64) model.Response {
	err := u.repository.FollowUser(userId, targetUserId)
	if err != nil {
		u.log.Errorf("repository.FollowUser: %v", err)

		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success follow user"),
	}
}

func (u *ProfileUsecase) UnfollowUser(userId, targetUserId int64) model.Response {
	err := u.repository.UnfollowUser(userId, targetUserId)
	if err != nil {
		u.log.Errorf("repository.UnfollowUser: %v", err)

		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success unfollow user"),
	}
}

func (u *ProfileUsecase) InsertUserWorkExperience(fileNames []string, props *model.WorkExperience) model.Response {
	var (
		err error
	)

	if len(fileNames) > 0 {
		objectPath := fmt.Sprintf("users/%d/work-experiences/files", props.UserId)

		defer func() {
			for _, fileName := range fileNames {
				filePath := fmt.Sprintf("./storage/temp/users/%d/files/%s", props.UserId, fileName)

				if err := u.fs.RemoveFile(filePath); err != nil {
					u.log.Errorf("fileSystem.RemoveFile: %v", err)
				}
			}
		}()

		urls, err := u.googleBucket.HandleObjectUploads(props.UserId, objectPath, fileNames...)
		if err != nil {
			u.log.Errorf("googleBucket.HandleObjectUploads: %v", err)

			return model.Response{
				Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
			}
		}

		props.FileURLs = urls
	}

	data, err := u.repository.InsertUserWorkExperience(props)
	if err != nil {
		// Delete uploaded objects
		errObjectDelete := u.googleBucket.HandleObjectDeletion(props.FileURLs...)
		if errObjectDelete != nil {
			u.log.Errorf("googleBucket.HandleObjectDeletion (user id: %d): %v", props.UserId, errObjectDelete)
		}

		u.log.Errorf("repository.InsertUserWorkExperience (user id: %d): %v", props.UserId, err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusCreated, "Success add user's work experience"),
		Data:   data,
	}
}

func (u *ProfileUsecase) InsertUserEducation(fileNames []string, props *model.Education) model.Response {
	var (
		err error
	)

	if len(fileNames) > 0 {
		defer func() {
			for _, fileName := range fileNames {
				filePath := fmt.Sprintf("./storage/temp/users/%d/files/%s", props.UserId, fileName)

				if err := u.fs.RemoveFile(filePath); err != nil {
					u.log.Errorf("fileSystem.RemoveFile: %v", err)
				}
			}
		}()

		objectPath := fmt.Sprintf("users/%d/educations/files", props.UserId)

		urls, err := u.googleBucket.HandleObjectUploads(props.UserId, objectPath, fileNames...)
		if err != nil {
			u.log.Errorf("googleBucket.HandleObjectUploads: %v", err)

			return model.Response{
				Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
			}
		}

		props.FileURLs = urls
	}

	data, err := u.repository.InsertUserEducation(props)
	if err != nil {
		// Delete uploaded objects
		errObjectDelete := u.googleBucket.HandleObjectDeletion(props.FileURLs...)
		if errObjectDelete != nil {
			u.log.Errorf("googleBucket.HandleObjectDeletion (user id: %d): %v", props.UserId, errObjectDelete)
		}

		u.log.Errorf("repository.InsertUserEducation (user id: %d): %v", props.UserId, err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusCreated, "Success add user's education"),
		Data:   data,
	}
}

func (u *ProfileUsecase) InsertUserProfile(imageFileNames []string, props *model.AddProfileRequest) model.Response {
	_, err := u.repository.GetUserById(props.UserId)
	if err != nil && err == sql.ErrNoRows {
		return model.Response{
			Status: libs.CustomResponse(http.StatusNotFound, "Data not found"),
		}
	}

	existingUser, err := u.repository.GetUserByEmail(props.Email)
	if err != nil && err != sql.ErrNoRows {
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	if existingUser.ID != props.UserId {
		return model.Response{
			Status: libs.CustomResponse(http.StatusBadRequest, "Email already exists"),
		}
	}

	if len(imageFileNames) > 0 {
		objectPath := fmt.Sprintf("users/%d/avatar", props.UserId)

		defer func() {
			filePath := fmt.Sprintf("./storage/temp/users/%d/files/%s", props.UserId, imageFileNames[0])

			if err := u.fs.RemoveFile(filePath); err != nil {
				u.log.Errorf("fileSystem.RemoveFile: %v", err)
			}
		}()

		urls, err := u.googleBucket.HandleObjectUploads(props.UserId, objectPath, imageFileNames[0])
		if err != nil {
			u.log.Errorf("googleBucket.HandleObjectUploads: %v", err)

			return model.Response{
				Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
			}
		}

		props.AvatarUrl = urls[0]
	}

	err = u.repository.InsertUserPersonalData(props)
	if err != nil {
		u.log.Errorf("repository.InsertUserPersonalData (user id: %d): %v", props.UserId, err)

		errObjectDelete := u.googleBucket.HandleObjectDeletion(props.AvatarUrl)
		if errObjectDelete != nil {
			u.log.Errorf("googleBucket.HandleObjectDeletion: %v", errObjectDelete)
		}

		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusCreated, "Success add user's personal data"),
		Data:   props,
	}
}
