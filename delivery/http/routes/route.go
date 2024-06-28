package routes

import (
	"database/sql"
	"profiln-be/delivery/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewRoute(app *gin.Engine, db *sql.DB, log *logrus.Logger) {
	app.Use(middleware.CORS())

	v1 := app.Group("/api/v1")

	NewAuthRoute(v1, db, log)
	NewHomepageRoute(v1, db, log)
	NewPostsRoute(v1, db, log)
	NewProfileRoute(v1, db, log)
	NewDataRoute(v1, db, log)
}
