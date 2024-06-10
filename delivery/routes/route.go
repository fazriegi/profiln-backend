package routes

import (
	"database/sql"
	"profiln-be/delivery/ws"
	ws_middleware "profiln-be/delivery/ws/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewRoute(app *gin.Engine, db *sql.DB, log *logrus.Logger) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	app.Use(cors.New(config))

	hub := ws.NewHub()
	go hub.Run()

	PostCommentsHandler := ws.NewPostCommentsHandler(hub, log)

	app.GET("/ws/posts/:postId/comments", ws_middleware.Authentication(), PostCommentsHandler.GetPostComments)

	v1 := app.Group("/api/v1")

	NewAuthRoute(v1, db, log)
	NewHomepageRoute(v1, db, log)
	NewPostsRoute(v1, db, log, hub)
	NewProfileRoute(v1, db, log)
	NewDataRoute(v1, db, log)
}
