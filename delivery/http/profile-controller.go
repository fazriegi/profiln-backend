package http

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"profiln-be/libs"
	"profiln-be/model"
	"profiln-be/package/profile"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type IProfileController interface {
	InsertUserAbout(ctx *gin.Context)
	UpdateProfile(ctx *gin.Context)
	UpdateAboutMe(ctx *gin.Context)
	UpdateUserCertificate(ctx *gin.Context)
	UpdateUserInformation(ctx *gin.Context)
	UpdateUserEducation(ctx *gin.Context)
	UpdateUserWorkExperience(ctx *gin.Context)
	AddUserOpenToWork(ctx *gin.Context)
	InsertUserCertificate(ctx *gin.Context)
	InsertUserSkills(ctx *gin.Context)
	GetUserProfile(ctx *gin.Context)
	GetUserWorkExperiences(ctx *gin.Context)
	GetUserEducations(ctx *gin.Context)
	GetUserCertificates(ctx *gin.Context)
	GetFollowedUsersByUser(ctx *gin.Context)
	GetUserBasicInformation(ctx *gin.Context)
	DeleteUserOpenToWork(ctx *gin.Context)
	DeleteUserWorkExperience(ctx *gin.Context)
	DeleteUserEducation(ctx *gin.Context)
	DeleteUserCertificate(ctx *gin.Context)
	FollowUser(ctx *gin.Context)
	UnfollowUser(ctx *gin.Context)
	InsertUserWorkExperience(ctx *gin.Context)
	InsertUserEducation(ctx *gin.Context)
}

type ProfileController struct {
	usecase profile.IProfileUsecase
}

func NewProfileController(usecase profile.IProfileUsecase) IProfileController {
	return &ProfileController{
		usecase,
	}
}

func (c *ProfileController) InsertUserSkills(ctx *gin.Context) {
	var (
		reqBody  model.SkillRequest
		response model.Response
	)

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Status = libs.CustomResponse(http.StatusBadRequest, "Error parsing request body")

		ctx.JSON(response.Status.Code, response)
		return
	}

	validationErr := libs.ValidateRequest(&reqBody) // validate reqBody struct
	// if there is an error
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		response.Status =
			libs.CustomResponse(http.StatusUnprocessableEntity, "Validation error")
		response.Data = errResponse

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.InsertUserSkill(&reqBody, userId)

	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) InsertUserCertificate(ctx *gin.Context) {
	var (
		response model.Response
		reqBody  model.Certificate
	)

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Error parsing request body")
		fmt.Println(err)
		ctx.JSON(response.Status.Code, response)
		return
	}

	validationErr := libs.ValidateRequest(reqBody) // validate reqBody struct
	// if there is an error
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		response.Status =
			libs.CustomResponse(http.StatusUnprocessableEntity, "Validation error")
		response.Data = errResponse

		ctx.JSON(response.Status.Code, response)
		return
	}

	reqBody.UserId = userId

	response = c.usecase.InsertUserCertificate(&reqBody)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) InsertUserAbout(ctx *gin.Context) {
	var (
		reqBody  model.UserDetailAboutRequest
		response model.Response
	)

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Status = libs.CustomResponse(http.StatusBadRequest, "Error parsing request body")

		ctx.JSON(response.Status.Code, response)
		return
	}

	validationErr := libs.ValidateRequest(&reqBody) // validate reqBody struct
	// if there is an error
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		response.Status =
			libs.CustomResponse(http.StatusUnprocessableEntity, "Validation error")
		response.Data = errResponse

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.InsertUserDetailAbout(&reqBody, userId)

	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) UpdateProfile(ctx *gin.Context) {
	var (
		response  model.Response
		reqBody   model.UpdateProfileRequest
		imageFile *multipart.FileHeader
	)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	// Get the first file
	files := ctx.MustGet("files").([]*multipart.FileHeader)
	if files != nil {
		imageFile = files[0]
	}

	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Error parsing request body")

		ctx.JSON(response.Status.Code, response)
		return
	}

	validationErr := libs.ValidateRequest(reqBody) // validate reqBody struct
	// if there is an error
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		response.Status =
			libs.CustomResponse(http.StatusUnprocessableEntity, "Validation error")
		response.Data = errResponse

		ctx.JSON(response.Status.Code, response)
		return
	}

	reqBody.UserId = userId

	response = c.usecase.UpdateProfile(imageFile, &reqBody)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) UpdateAboutMe(ctx *gin.Context) {
	var (
		response model.Response
		reqBody  model.UserDetailAboutRequest
	)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Error parsing request body")

		ctx.JSON(response.Status.Code, response)
		return
	}

	validationErr := libs.ValidateRequest(reqBody) // validate reqBody struct
	// if there is an error
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		response.Status =
			libs.CustomResponse(http.StatusUnprocessableEntity, "Validation error")
		response.Data = errResponse

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.UpdateAboutMe(userId, reqBody.About)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) UpdateUserCertificate(ctx *gin.Context) {
	var (
		response model.Response
		reqBody  model.Certificate
	)

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	certificateId, err := strconv.ParseInt(ctx.Param("certificateId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Error parsing request body")
		fmt.Println(err)
		ctx.JSON(response.Status.Code, response)
		return
	}

	validationErr := libs.ValidateRequest(reqBody) // validate reqBody struct
	// if there is an error
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		response.Status =
			libs.CustomResponse(http.StatusUnprocessableEntity, "Validation error")
		response.Data = errResponse

		ctx.JSON(response.Status.Code, response)
		return
	}

	reqBody.ID = certificateId

	response = c.usecase.UpdateUserCertificate(userId, &reqBody)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) UpdateUserInformation(ctx *gin.Context) {
	var (
		response model.Response
		reqBody  model.UpdateUserInformation
	)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Error parsing request body")

		ctx.JSON(response.Status.Code, response)
		return
	}

	validationErr := libs.ValidateRequest(reqBody) // validate reqBody struct
	// if there is an error
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		response.Status =
			libs.CustomResponse(http.StatusUnprocessableEntity, "Validation error")
		response.Data = errResponse

		ctx.JSON(response.Status.Code, response)
		return
	}

	reqBody.UserId = userId

	response = c.usecase.UpdateUserInformation(&reqBody)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) UpdateUserEducation(ctx *gin.Context) {
	var (
		response model.Response
		reqBody  model.Education
	)
	files := ctx.MustGet("files").([]*multipart.FileHeader)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	educationId, err := strconv.ParseInt(ctx.Param("educationId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Error parsing request body")

		ctx.JSON(response.Status.Code, response)
		return
	}

	validationErr := libs.ValidateRequest(reqBody) // validate reqBody struct
	// if there is an error
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		response.Status =
			libs.CustomResponse(http.StatusUnprocessableEntity, "Validation error")
		response.Data = errResponse

		ctx.JSON(response.Status.Code, response)
		return
	}

	reqBody.UserId = userId
	reqBody.ID = educationId

	response = c.usecase.UpdateUserEducation(files, &reqBody)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) UpdateUserWorkExperience(ctx *gin.Context) {
	var (
		response model.Response
		reqBody  model.WorkExperience
	)
	files := ctx.MustGet("files").([]*multipart.FileHeader)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	workExperienceId, err := strconv.ParseInt(ctx.Param("workExperienceId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Error parsing request body")

		ctx.JSON(response.Status.Code, response)
		return
	}

	validationErr := libs.ValidateRequest(reqBody) // validate reqBody struct
	// if there is an error
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		response.Status =
			libs.CustomResponse(http.StatusUnprocessableEntity, "Validation error")
		response.Data = errResponse

		ctx.JSON(response.Status.Code, response)
		return
	}

	reqBody.UserId = userId
	reqBody.ID = workExperienceId

	response = c.usecase.UpdateUserWorkExperience(files, &reqBody)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) GetUserProfile(ctx *gin.Context) {
	var (
		response model.Response
	)

	userId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.GetUserProfile(userId)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) GetUserWorkExperiences(ctx *gin.Context) {
	var response model.Response

	userId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	if page <= 0 || limit <= 0 {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	pagination := model.PaginationRequest{
		Page:  page,
		Limit: limit,
	}
	response = c.usecase.GetWorkExperiencesByUserId(userId, pagination)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) GetUserEducations(ctx *gin.Context) {
	var response model.Response

	userId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	if page <= 0 || limit <= 0 {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	pagination := model.PaginationRequest{
		Page:  page,
		Limit: limit,
	}
	response = c.usecase.GetEducationsByUserId(userId, pagination)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) GetUserCertificates(ctx *gin.Context) {
	var response model.Response

	userId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	if page <= 0 || limit <= 0 {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	pagination := model.PaginationRequest{
		Page:  page,
		Limit: limit,
	}
	response = c.usecase.GetCertificatesByUserId(userId, pagination)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) GetFollowedUsersByUser(ctx *gin.Context) {
	var response model.Response

	userId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	if page <= 0 || limit <= 0 {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	pagination := model.PaginationRequest{
		Page:  page,
		Limit: limit,
	}
	response = c.usecase.GetFollowedUsersByUserId(userId, pagination)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) GetUserBasicInformation(ctx *gin.Context) {
	var (
		response model.Response
	)

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	response = c.usecase.GetUserBasicInformation(userId)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) AddUserOpenToWork(ctx *gin.Context) {
	var (
		response model.Response
		reqBody  model.OpenToWork
	)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Error parsing request body")

		ctx.JSON(response.Status.Code, response)
		return
	}

	validationErr := libs.ValidateRequest(reqBody) // validate reqBody struct
	// if there is an error
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		response.Status =
			libs.CustomResponse(http.StatusUnprocessableEntity, "Validation error")
		response.Data = errResponse

		ctx.JSON(response.Status.Code, response)
		return
	}

	reqBody.UserId = userId

	response = c.usecase.AddUserOpenToWork(&reqBody)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) DeleteUserOpenToWork(ctx *gin.Context) {
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	response := c.usecase.DeleteUserOpenToWork(userId)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) DeleteUserWorkExperience(ctx *gin.Context) {
	var response model.Response
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	workExperienceId, err := strconv.ParseInt(ctx.Param("workExperienceId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.DeleteUserWorkExperienceById(userId, workExperienceId)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) DeleteUserEducation(ctx *gin.Context) {
	var response model.Response
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	educationId, err := strconv.ParseInt(ctx.Param("educationId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.DeleteUserEducationById(userId, educationId)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) DeleteUserCertificate(ctx *gin.Context) {
	var response model.Response
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	certificateId, err := strconv.ParseInt(ctx.Param("certificateId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.DeleteUserCertificateById(userId, certificateId)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) FollowUser(ctx *gin.Context) {
	var response model.Response
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	targetUserId, err := strconv.ParseInt(ctx.Param("targetUserId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	if userId == targetUserId {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Can't follow yourself")

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.FollowUser(userId, targetUserId)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) UnfollowUser(ctx *gin.Context) {
	var response model.Response
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	targetUserId, err := strconv.ParseInt(ctx.Param("targetUserId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	if userId == targetUserId {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Can't unfollow yourself")

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.UnfollowUser(userId, targetUserId)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) InsertUserWorkExperience(ctx *gin.Context) {
	var (
		response model.Response
		reqBody  model.WorkExperience
	)
	files := ctx.MustGet("files").([]*multipart.FileHeader)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Error parsing request body")

		ctx.JSON(response.Status.Code, response)
		return
	}

	validationErr := libs.ValidateRequest(reqBody) // validate reqBody struct
	// if there is an error
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		response.Status =
			libs.CustomResponse(http.StatusUnprocessableEntity, "Validation error")
		response.Data = errResponse

		ctx.JSON(response.Status.Code, response)
		return
	}

	reqBody.UserId = userId

	response = c.usecase.InsertUserWorkExperience(files, &reqBody)
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) InsertUserEducation(ctx *gin.Context) {
	var (
		response model.Response
		reqBody  model.Education
	)
	files := ctx.MustGet("files").([]*multipart.FileHeader)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Error parsing request body")

		ctx.JSON(response.Status.Code, response)
		return
	}

	validationErr := libs.ValidateRequest(reqBody) // validate reqBody struct
	// if there is an error
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		response.Status =
			libs.CustomResponse(http.StatusUnprocessableEntity, "Validation error")
		response.Data = errResponse

		ctx.JSON(response.Status.Code, response)
		return
	}

	reqBody.UserId = userId

	response = c.usecase.InsertUserEducation(files, &reqBody)
	ctx.JSON(response.Status.Code, response)
}
