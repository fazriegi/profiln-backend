package routes

import (
	"database/sql"
	"profiln-be/delivery/http"
	"profiln-be/delivery/http/middleware"
	"profiln-be/package/homepage"
	repository "profiln-be/package/homepage/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewHomepageRoute(app *gin.RouterGroup, db *sql.DB, log *logrus.Logger) {
	repository := repository.NewHomepageRepository(db)
	usecase := homepage.NewHomepageUsecase(repository, log)
	controller := http.NewHomepageController(usecase)

	app.Use(middleware.Authentication())
	app.GET("/posts", controller.ListPosts)
	app.GET("/users/me/follow-recommendations", controller.ListFollowsRecommendation)
}
