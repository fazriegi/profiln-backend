package routes

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewRoute(app *gin.Engine, db *sql.DB, log *logrus.Logger) {
	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	app.Use(cors.New(corsConfig))

	v1 := app.Group("/api/v1")

	NewAuthRoute(v1, db, log)
	NewHomepageRoute(v1, db, log)
	NewPostsRoute(v1, db, log)
	NewProfileRoute(v1, db, log)
	NewDataRoute(v1, db, log)
}
