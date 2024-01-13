package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type wsServer struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	register  chan *websocket.Conn
}

var wss = wsServer{
	clients:   make(map[*websocket.Conn]bool),
	broadcast: make(chan []byte),
	register:  make(chan *websocket.Conn),
}

func (wss *wsServer) run() {
	for {
		select {
		case c := <-wss.register:
			wss.clients[c] = true
		case msg := <-wss.broadcast:
			for c := range wss.clients {
				err := c.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					delete(wss.clients, c)
					c.Close()
				}
			}
		}
	}
}

func serveWs(c *gin.Context) {
	w, r := c.Writer, c.Request
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	wss.register <- conn
}
