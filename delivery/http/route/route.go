package route

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRoute(app *gin.Engine, db *sql.DB) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}

	app.Use(cors.New(config))

	v1 := app.Group("/api/v1")

	NewUserRoute(v1, db)
	NewEmailRoute(v1)
}
