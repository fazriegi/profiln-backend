package routes

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewRoute(app *gin.Engine, db *sql.DB, log *logrus.Logger) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	app.Use(cors.New(config))

	v1 := app.Group("/api/v1")

	NewAuthRoute(v1, db, log)
	NewHomepageRoute(v1, db, log)
	NewPostsRoute(v1, db, log)
	NewProfileRoute(v1, db, log)
	NewDataRoute(v1, db, log)
}
