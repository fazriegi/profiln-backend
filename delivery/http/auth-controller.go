package http

import (
	"net/http"

	"profiln-be/libs"
	"profiln-be/model"
	"profiln-be/package/auth"

	"github.com/gin-gonic/gin"
)

type IAuthController interface {
	Login(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
	Register(ctx *gin.Context)
	VerifiedEmail(ctx *gin.Context)
	SendResetPasswordEmail(ctx *gin.Context)
	SendOTPEmail(ctx *gin.Context)
}

type AuthController struct {
	usecase auth.IAuthUsecase
}

func NewAuthController(usecase auth.IAuthUsecase) IAuthController {
	return &AuthController{
		usecase,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var (
		reqBody  model.LoginRequest
		response model.Response
	)

	loginType := ctx.Query("type")

	if loginType != "sso" && loginType != "app" {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Login type is not valid")

		ctx.JSON(response.Status.Code, response)
		return
	}

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
	response = c.usecase.Login(loginType, &reqBody)

	ctx.JSON(response.Status.Code, response)
}

func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var (
		reqBody  model.ResetPasswordRequest
		response model.Response
	)

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

	response = c.usecase.ResetPassword(&reqBody)
	ctx.JSON(response.Status.Code, response)
}

func (c *AuthController) Register(ctx *gin.Context) {
	var (
		reqBody  model.RegisterRequest
		response model.Response
	)

	err := ctx.ShouldBind(&reqBody)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Error parsing request body")

		ctx.JSON(response.Status.Code, response)
		return
	}

	oauth := ctx.Query("oauth")
	if oauth != "true" && oauth != "false" {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	validationErr := libs.ValidateRequest(reqBody)

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

	if oauth == "false" && reqBody.Password == "" {
		response.Status =
			libs.CustomResponse(http.StatusUnprocessableEntity, "Validation error")
		response.Data = map[string]string{
			"errors": "password can't be empty",
		}
		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.Register(&reqBody, oauth)

	ctx.JSON(response.Status.Code, response)
}

func (c *AuthController) VerifiedEmail(ctx *gin.Context) {
	var (
		reqBody  model.VerifiedEmailOTPRequest
		response model.Response
	)

	err := ctx.ShouldBind(&reqBody)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Error parsing request body")

		ctx.JSON(response.Status.Code, response)
		return
	}

	validationErr := libs.ValidateRequest(reqBody)

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

	response = c.usecase.UpdateVerifiedEmail(&reqBody)

	ctx.JSON(response.Status.Code, response)
}

func (c *AuthController) SendResetPasswordEmail(ctx *gin.Context) {
	var (
		reqBody  model.ResetPasswordEmailRequest
		response model.Response
	)

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

	response = c.usecase.SendResetPasswordEmail(&reqBody)

	ctx.JSON(response.Status.Code, response)
}

func (c *AuthController) SendOTPEmail(ctx *gin.Context) {
	var (
		reqBody  model.OTPEmailRequest
		response model.Response
	)

	if err := ctx.ShouldBind(&reqBody); err != nil {

		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Error parsing request body")

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

	response = c.usecase.SendOTPEmail(&reqBody)

	ctx.JSON(response.Status.Code, response)
}
