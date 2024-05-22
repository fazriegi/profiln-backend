package http

import (
	"net/http"
	"profiln-be/libs"
	"profiln-be/model"
	"profiln-be/package/profile"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type IProfileController interface {
	InsertUserAbout(ctx *gin.Context)
	GetSkills(ctx *gin.Context)
	UpdateProfile(ctx *gin.Context)
	UpdateAboutMe(ctx *gin.Context)
}

type ProfileController struct {
	usecase profile.IProfileUsecase
}

func NewProfileController(usecase profile.IProfileUsecase) IProfileController {
	return &ProfileController{
		usecase,
	}
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

func (c *ProfileController) GetSkills(ctx *gin.Context) {
	response := c.usecase.GetSkills()
	ctx.JSON(response.Status.Code, response)
}

func (c *ProfileController) UpdateProfile(ctx *gin.Context) {
	var (
		response model.Response
		reqBody  model.UpdateProfileRequest
	)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	imageFile, _ := ctx.FormFile("file")

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
