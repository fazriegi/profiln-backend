package route

import (
	"database/sql"

	"profiln-be/delivery/http"
	"profiln-be/usecase"

	"github.com/gin-gonic/gin"
)

func NewUserRoute(app *gin.RouterGroup, db *sql.DB) {
	userUsecase := usecase.NewUserUsecase(db)
	userController := http.NewUserController(userUsecase)

	app.POST("/login", userController.Login)
	app.POST("/reset-password", userController.ResetPassword)
}
