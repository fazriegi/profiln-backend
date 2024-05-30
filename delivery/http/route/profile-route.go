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

	profile := app.Group("profiles")
	profile.Use(middleware.Authentication())
	profile.POST("/user/certificate", controller.InsertCertificate)
	profile.POST("/user/skill", controller.InsertUserSkills)
	// profile.GET("/user/about", controller.GetUserAbout)
	// profile.GET("/user/certificate", controller.GetUserCertificates)
	// profile.GET("/user/skill", controller.GetUserSkillsLocationPortofolio)

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
	me.GET("/", controller.GetUserBasicInformation)

	users := app.Group("users")
	users.GET("/:userId/profile", controller.GetUserProfile)
	users.GET("/:userId/work-experiences", controller.GetUserWorkExperiences)
	users.GET("/:userId/educations", controller.GetUserEducations)
	users.GET("/:userId/certificates", controller.GetUserCertificates)
	users.GET("/:userId/followings", controller.GetFollowedUsersByUser)
}
