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
	googleBucket := libs.NewGoogleBucket(fileSystem, log)
	repository := repository.NewPostsRepository(db)
	usecase := posts.NewPostsUsecase(repository, log, googleBucket)
	controller := http.NewPostsController(usecase)

	app.Use(middleware.Authentication())

	posts := app.Group("posts")
	posts.POST("/:postId/report", controller.ReportPost)
	posts.GET("/:postId", controller.GetDetailPost)
	posts.GET("/:postId/comments", controller.GetPostComments)
	posts.GET("/:postId/comments/:postCommentId/replies", controller.GetPostCommentReplies)
	posts.PATCH("/:postId/like", controller.LikePost)

	myPosts := app.Group("users/me/posts")
	myPosts.GET("/", controller.ListNewestPostsByUserId)
	myPosts.GET("/like", controller.ListLikedPostsByUserId)
	myPosts.GET("/repost", controller.ListRepostedPostsByUserId)
	myPosts.POST("/", middleware.ValidateFileUpload(int64(twoMegaBytes), 10, imageFormats), controller.InsertPost)
	myPosts.PATCH("/:postId", middleware.ValidateFileUpload(int64(twoMegaBytes), 10, imageFormats), controller.UpdatePost)
	myPosts.DELETE("/:postId", controller.DeletePost)
}
