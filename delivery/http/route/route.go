package route

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewRoute(app *gin.Engine, db *sql.DB, log *logrus.Logger) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	app.Use(cors.New(config))

	v1 := app.Group("/api/v1")

	NewAuthRoute(v1, db, log)
	NewProfileRoute(v1, db, log)
}
