package route

import (
	"database/sql"
	"profiln-be/delivery/http"
	"profiln-be/usecase"

	"github.com/gin-gonic/gin"
)

func NewEmailRoute(app *gin.RouterGroup, db *sql.DB) {
	emailUsecase := usecase.NewEmailUsecase(db)
	emailController := http.NewEmailController(emailUsecase)

	email := app.Group("/email")
	email.POST("/reset-password", emailController.SendResetPasswordMail)
}
