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
	twoMegaBytes := 2 << 20

	repository := repository.NewProfileRepository(db)
	usecase := profile.NewProfileUsecase(repository, log)
	controller := http.NewProfileController(usecase)

	profile := app.Group("profiles")
	profile.Use(middleware.Authentication())
	profile.POST("/users/about", controller.InsertUserAbout)
	profile.GET("/skills", controller.GetSkills)
	profile.PUT("/my-profile", middleware.MaxReqSizeAllowed(int64(twoMegaBytes)), controller.UpdateProfile)
	profile.PUT("/about", controller.UpdateAboutMe)
}
