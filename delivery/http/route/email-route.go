package route

import (
	"profiln-be/delivery/http"
	"profiln-be/usecase"

	"github.com/gin-gonic/gin"
)

func NewEmailRoute(app *gin.RouterGroup) {
	emailUsecase := usecase.NewEmailUsecase()
	emailController := http.NewEmailController(emailUsecase)

	email := app.Group("/email")
	email.POST("/reset-password", emailController.SendResetPasswordMail)
}
