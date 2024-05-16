package route

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

	homepage := app.Group("homepage")
	homepage.Use(middleware.Authentication())
	homepage.GET("/posts", controller.ListPosts)
	homepage.GET("/follows-recommendation", controller.ListFollowsRecommendation)
}
