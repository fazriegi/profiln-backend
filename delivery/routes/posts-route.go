package routes

import (
	"database/sql"
	"profiln-be/delivery/http"
	"profiln-be/delivery/http/middleware"
	"profiln-be/delivery/ws"
	"profiln-be/libs"
	"profiln-be/package/posts"
	repository "profiln-be/package/posts/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewPostsRoute(app *gin.RouterGroup, db *sql.DB, log *logrus.Logger, wsHub *ws.Hub) {
	twoMegaBytes := 2 * 1024 * 1024
	imageFormats := []string{".png", ".jpg"}

	fileSystem := libs.NewFileSystem()
	googleBucket := libs.NewGoogleBucket(log)
	repository := repository.NewPostsRepository(db)
	usecase := posts.NewPostsUsecase(repository, log, googleBucket, fileSystem)
	controller := http.NewPostsController(usecase, wsHub)

	app.Use(middleware.Authentication())

	app.GET("/users/:userId/posts", controller.ListNewestPostsByTargetUser)
	app.GET("/users/:userId/posts/like", controller.ListLikedPostsByTargetUser)
	app.GET("/users/:userId/posts/repost", controller.ListRepostedPostsByTargetUser)

	posts := app.Group("posts")
	posts.POST("/:postId/report", controller.ReportPost)
	posts.GET("/:postId", controller.GetDetailPost)
	posts.GET("/:postId/comments", controller.GetPostComments)
	posts.GET("/:postId/comments/:postCommentId/replies", controller.GetPostCommentReplies)
	posts.POST("/:postId/like", controller.LikePost)
	posts.POST("/:postId/unlike", controller.UnlikePost)
	posts.POST("/:postId/repost", controller.RepostPost)
	posts.POST("/:postId/unrepost", controller.UnrepostPost)
	posts.POST("/:postId/comments", middleware.ValidateFileUpload(int64(twoMegaBytes), 1, imageFormats, fileSystem, log), controller.InsertPostComment)

	myPosts := app.Group("users/me/posts")
	myPosts.POST("/", controller.InsertPost)
	myPosts.PATCH("/:postId", controller.UpdatePost)
	myPosts.DELETE("/:postId", controller.DeletePost)
	myPosts.POST("/:postId/upload", middleware.ValidateFileUpload(int64(twoMegaBytes), 10, imageFormats, fileSystem, log), controller.UploadFileForInsertPost)
	myPosts.PUT("/:postId/upload", middleware.ValidateFileUpload(int64(twoMegaBytes), 10, imageFormats, fileSystem, log), controller.UploadFileForUpdatePost)
}
