package http

import (
	"net/http"

	"profiln-be/libs"
	"profiln-be/model"
	"profiln-be/usecase"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUsecase *usecase.UserUsecase
}

func NewUserController(userUsecase *usecase.UserUsecase) *UserController {
	return &UserController{
		UserUsecase: userUsecase,
	}
}

func (c *UserController) Login(ctx *gin.Context) {
	var (
		reqBody  model.UserLoginRequest
		response model.Response
	)

	if err := ctx.ShouldBind(&reqBody); err != nil {
		status := libs.CustomResponse(http.StatusBadRequest, "Error parsing request body")
		response.Status = status

		ctx.JSON(status.Code, response)
		return
	}

	validationErr := libs.ValidateRequest(reqBody) // validate reqBody struct
	// if there is an error
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		status := libs.CustomResponse(http.StatusUnprocessableEntity, "Validation error")
		response.Status = status
		response.Data = errResponse

		ctx.JSON(status.Code, response)
		return
	}

	response = c.UserUsecase.Login(&reqBody)

	ctx.JSON(response.Status.Code, response)
}

func (c *UserController) ResetPassword(ctx *gin.Context) {
	var (
		reqBody  model.UserResetPasswordRequest
		response model.Response
	)

	if err := ctx.ShouldBind(&reqBody); err != nil {
		status := libs.CustomResponse(http.StatusBadRequest, "Error parsing request body")
		response.Status = status

		ctx.JSON(status.Code, response)
		return
	}

	validationErr := libs.ValidateRequest(reqBody) // validate reqBody struct
	// if there is an error
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		status := libs.CustomResponse(http.StatusUnprocessableEntity, "Validation error")
		response.Status = status
		response.Data = errResponse

		ctx.JSON(status.Code, response)
		return
	}

	response = c.UserUsecase.ResetPassword(&reqBody)

	ctx.JSON(response.Status.Code, response)
}
