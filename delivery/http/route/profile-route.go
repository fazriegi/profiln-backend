package route

import (
	"database/sql"
	"profiln-be/delivery/http"
	"profiln-be/delivery/http/middleware"
	"profiln-be/package/profile"
	repository "profiln-be/package/profile/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewProfileRoute(app *gin.RouterGroup, db *sql.DB, log *logrus.Logger) {
	repository := repository.NewProfileRepository(db)
	usecase := profile.NewProfileUsecase(repository, log)
	controller := http.NewProfileController(usecase)

	profile := app.Group("profile")
	profile.Use(middleware.Authentication())
	profile.POST("/user/about", controller.InsertUserAbout)
}
