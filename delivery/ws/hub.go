package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	subscribe  map[int64]map[*Client]bool
}

type Message struct {
	PostId int64 `json:"post_id"`
	Data   any   `json:"data"`
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	postId int64
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		subscribe:  make(map[int64]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			if _, ok := h.subscribe[client.postId]; !ok {
				h.subscribe[client.postId] = make(map[*Client]bool)
			}
			h.subscribe[client.postId][client] = true

		case client := <-h.unregister:
			delete(h.clients, client)
			delete(h.subscribe[client.postId], client)

		case message := <-h.broadcast:
			if clients, ok := h.subscribe[message.PostId]; ok {
				for client := range clients {
					client.conn.WriteJSON(message)
				}
			}

		}
	}
}

func (h *Hub) Broadcast(message Message) {
	h.broadcast <- message
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
