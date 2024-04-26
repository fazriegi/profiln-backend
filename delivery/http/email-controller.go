package http

import (
	"net/http"
	"profiln-be/libs"
	"profiln-be/model"
	"profiln-be/usecase"

	"github.com/gin-gonic/gin"
)

type EmailController struct {
	EmailUsecase *usecase.EmailUsecase
}

func NewEmailController(EmailUsecase *usecase.EmailUsecase) *EmailController {
	return &EmailController{
		EmailUsecase: EmailUsecase,
	}
}

func (c *EmailController) SendResetPasswordMail(ctx *gin.Context) {
	var (
		reqBody  model.SendResetPassEmailRequest
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

	response = c.EmailUsecase.SendResetPasswordMail(&reqBody)

	ctx.JSON(response.Status.Code, response)
}
