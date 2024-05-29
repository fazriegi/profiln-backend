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
	twoMegaBytes := 2 * 1024 * 1024
	imageFormats := []string{".png", ".jpg"}
	imageAndDocumentFormats := append(imageFormats, ".pdf", ".doc", ".docx")

	fileSystem := libs.NewFileSystem()
	googleBucket := libs.NewGoogleBucket(fileSystem, log)
	repository := repository.NewProfileRepository(db)
	usecase := profile.NewProfileUsecase(repository, log, googleBucket)
	controller := http.NewProfileController(usecase)

	app.Use(middleware.Authentication())
	app.GET("/skills", controller.GetSkills)

	me := app.Group("users/me")
	me.POST("/about", controller.InsertUserAbout)
	me.PUT("/profile", middleware.ValidateFileUpload(int64(twoMegaBytes), 1, imageFormats), controller.UpdateProfile)
	me.PUT("/about", controller.UpdateAboutMe)
	me.PUT("/certificates/:certificateId", controller.UpdateUserCertificate)
	me.PUT("/information", controller.UpdateUserInformation)
	me.PUT("/educations/:educationId", middleware.ValidateFileUpload(int64(twoMegaBytes), 3, imageAndDocumentFormats), controller.UpdateUserEducation)
	me.PUT("/work-experiences/:workExperienceId", middleware.ValidateFileUpload(int64(twoMegaBytes), 3, imageAndDocumentFormats), controller.UpdateUserWorkExperience)
}