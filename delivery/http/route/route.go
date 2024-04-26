package route

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func NewRoute(app *gin.Engine, db *sql.DB) {
	v1 := app.Group("/api/v1")

	NewUserRoute(v1, db)
	NewEmailRoute(v1)
}
