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
	googleBucket := libs.NewGoogleBucket(log)
	repository := repository.NewProfileRepository(db)
	usecase := profile.NewProfileUsecase(repository, log, googleBucket, fileSystem)
	controller := http.NewProfileController(usecase)

	profile := app.Group("profiles")
	profile.Use(middleware.Authentication())
	profile.POST("/user/skill", controller.InsertUserSkills)

	app.Use(middleware.Authentication())

	me := app.Group("users/me")
	me.POST("/about", controller.InsertUserAbout)
	me.PUT("/profile", middleware.ValidateFileUpload(int64(twoMegaBytes), 1, imageFormats, fileSystem, log), controller.UpdateProfile)
	me.PUT("/about", controller.UpdateAboutMe)
	me.PUT("/certificates/:certificateId", controller.UpdateUserCertificate)
	me.PUT("/information", controller.UpdateUserInformation)
	me.PUT("/educations/:educationId", middleware.ValidateFileUpload(int64(twoMegaBytes), 3, imageAndDocumentFormats, fileSystem, log), controller.UpdateUserEducation)
	me.PUT("/work-experiences/:workExperienceId", middleware.ValidateFileUpload(int64(twoMegaBytes), 3, imageAndDocumentFormats, fileSystem, log), controller.UpdateUserWorkExperience)
	me.POST("/open-to-work", controller.AddUserOpenToWork)
	me.DELETE("/open-to-work", controller.DeleteUserOpenToWork)
	me.DELETE("/work-experiences/:workExperienceId", controller.DeleteUserWorkExperience)
	me.DELETE("/educations/:educationId", controller.DeleteUserEducation)
	me.DELETE("/certificates/:certificateId", controller.DeleteUserCertificate)
	me.POST("/work-experiences", middleware.ValidateFileUpload(int64(twoMegaBytes), 3, imageAndDocumentFormats, fileSystem, log), controller.InsertUserWorkExperience)
	me.POST("/educations", middleware.ValidateFileUpload(int64(twoMegaBytes), 3, imageAndDocumentFormats, fileSystem, log), controller.InsertUserEducation)
	me.POST("/certificates", controller.InsertUserCertificate)
	me.GET("/", controller.GetUserBasicInformation)

	users := app.Group("users")
	users.GET("/:userId/profile", controller.GetUserProfile)
	users.GET("/:userId/work-experiences", controller.GetUserWorkExperiences)
	users.GET("/:userId/educations", controller.GetUserEducations)
	users.GET("/:userId/certificates", controller.GetUserCertificates)
	users.GET("/:userId/followings", controller.GetFollowedUsersByUser)
	users.POST("/:targetUserId/follow", controller.FollowUser)
	users.DELETE("/:targetUserId/follow", controller.UnfollowUser)
}
