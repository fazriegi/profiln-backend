package route

import (
	"database/sql"
	"profiln-be/delivery/http"
	"profiln-be/delivery/http/middleware"
	"profiln-be/libs"
	"profiln-be/package/profile"
	repository "profiln-be/package/profile/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewProfileRoute(app *gin.RouterGroup, db *sql.DB, log *logrus.Logger) {
	twoMegaBytes := 2 << 20

	fileSystem := libs.NewFileSystem()
	googleBucket := libs.NewGoogleBucket(fileSystem, log)
	repository := repository.NewProfileRepository(db)
	usecase := profile.NewProfileUsecase(repository, log, googleBucket)
	controller := http.NewProfileController(usecase)

	profile := app.Group("profiles")
	profile.Use(middleware.Authentication())
	profile.POST("/users/about", controller.InsertUserAbout)
	profile.GET("/skills", controller.GetSkills)
	profile.PUT("/my-profile", middleware.MaxReqSizeAllowed(int64(twoMegaBytes)), controller.UpdateProfile)
	profile.PUT("/about", controller.UpdateAboutMe)
	profile.PUT("/certificates/:certificateId", controller.UpdateUserCertificate)
	profile.PUT("/my-information", controller.UpdateUserInformation)
}
