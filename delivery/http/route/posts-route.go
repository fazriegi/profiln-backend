package route

import (
	"database/sql"
	"profiln-be/delivery/http"
	"profiln-be/delivery/http/middleware"
	"profiln-be/package/posts"
	repository "profiln-be/package/posts/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewPostsRoute(app *gin.RouterGroup, db *sql.DB, log *logrus.Logger) {
	repository := repository.NewPostsRepository(db)
	usecase := posts.NewPostsUsecase(repository, log)
	controller := http.NewPostsController(usecase)

	posts := app.Group("posts")
	posts.Use(middleware.Authentication())
	posts.POST("/:postId/report", controller.ReportPost)
	posts.GET("/:postId", controller.GetDetailPost)
	posts.GET("/:postId/comments", controller.GetPostComments)
	posts.GET("/:postId/comments/:postCommentId/replies", controller.GetPostCommentReplies)
	posts.POST("/:postId/like", controller.UpdatePostLikeCount)

	users := app.Group("users")
	users.Use(middleware.Authentication())
	users.GET("/me/posts", controller.ListNewestPostsByUserId)
	users.GET("/me/posts/liked", controller.ListLikedPostsByUserId)
}
