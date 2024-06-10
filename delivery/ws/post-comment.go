package ws

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type IPostCommentsHandler interface {
	GetPostComments(ctx *gin.Context)
}

type PostCommentsHandler struct {
	hub *Hub
	log *logrus.Logger
}

func NewPostCommentsHandler(hub *Hub, log *logrus.Logger) IPostCommentsHandler {
	return &PostCommentsHandler{
		hub,
		log,
	}
}

func (h *PostCommentsHandler) GetPostComments(ctx *gin.Context) {
	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid request param")
		return
	}

	conn, err := Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		h.log.Errorf("ws.Upgrader.Upgrade: %v", err)
		ctx.JSON(http.StatusInternalServerError, "Unexpected error occured")
		return
	}
	defer conn.Close()

	client := &Client{hub: h.hub, conn: conn, postId: postId}
	client.hub.register <- client

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			h.log.Errorf("ws.conn.ReadMessage(): %v", err)
			client.hub.unregister <- client
			ctx.JSON(http.StatusInternalServerError, "Unexpected error occured")
			return
		}
	}
}
