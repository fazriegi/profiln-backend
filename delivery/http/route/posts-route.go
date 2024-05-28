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

	posts := app.Group("posts")
	posts.Use(middleware.Authentication())
	posts.POST("/:postId/report", controller.ReportPost)
	posts.GET("/:postId", controller.GetDetailPost)
	posts.GET("/:postId/comments", controller.GetPostComments)
	posts.GET("/:postId/comments/:postCommentId/replies", controller.GetPostCommentReplies)
	posts.POST("/:postId/like", controller.UpdatePostLikeCount)

	myPosts := app.Group("users/me/posts")
	myPosts.Use(middleware.Authentication())
	myPosts.GET("/", controller.ListNewestPostsByUserId)
	myPosts.GET("/liked", controller.ListLikedPostsByUserId)
	myPosts.GET("/reposted", controller.ListRepostedPostsByUserId)
	myPosts.POST("/", middleware.ValidateFileUpload(int64(twoMegaBytes), 1, imageFormats), controller.InsertPost)
}
