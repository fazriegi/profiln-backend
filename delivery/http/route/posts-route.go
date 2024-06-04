package route

import (
	"database/sql"
	"profiln-be/delivery/http"
	"profiln-be/delivery/http/middleware"
	"profiln-be/libs"
	"profiln-be/package/posts"
	repository "profiln-be/package/posts/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewPostsRoute(app *gin.RouterGroup, db *sql.DB, log *logrus.Logger) {
	twoMegaBytes := 2 * 1024 * 1024
	imageFormats := []string{".png", ".jpg"}

	fileSystem := libs.NewFileSystem()
	googleBucket := libs.NewGoogleBucket(log)
	repository := repository.NewPostsRepository(db)
	usecase := posts.NewPostsUsecase(repository, log, googleBucket, fileSystem)
	controller := http.NewPostsController(usecase)

	app.Use(middleware.Authentication())

	posts := app.Group("posts")
	posts.POST("/:postId/report", controller.ReportPost)
	posts.GET("/:postId", controller.GetDetailPost)
	posts.GET("/:postId/comments", controller.GetPostComments)
	posts.GET("/:postId/comments/:postCommentId/replies", controller.GetPostCommentReplies)
	posts.POST("/:postId/like", controller.LikePost)
	posts.POST("/:postId/unlike", controller.UnlikePost)
	posts.POST("/:postId/repost", controller.RepostPost)
	posts.POST("/:postId/unrepost", controller.UnrepostPost)

	myPosts := app.Group("users/me/posts")
	myPosts.GET("/", controller.ListNewestPostsByUserId)
	myPosts.GET("/like", controller.ListLikedPostsByUserId)
	myPosts.GET("/repost", controller.ListRepostedPostsByUserId)
	myPosts.POST("/", controller.InsertPost)
	myPosts.PATCH("/:postId", controller.UpdatePost)
	myPosts.DELETE("/:postId", controller.DeletePost)
	myPosts.POST("/:postId/upload", middleware.ValidateFileUpload(int64(twoMegaBytes), 10, imageFormats, fileSystem, log), controller.UploadFileForInsertPost)
	myPosts.PUT("/:postId/upload", middleware.ValidateFileUpload(int64(twoMegaBytes), 10, imageFormats, fileSystem, log), controller.UploadFileForUpdatePost)
}
