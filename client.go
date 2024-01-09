package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel for outbound message from server.
	// For the counter app, the message is the value of the counter.
	send chan int
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// log only if error is unexpected
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Println("read:", err)
			}
			break
		}

		// the only message is to increment
		_ = message

		c.hub.broadcast <- true
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		message, ok := <-c.send
		if !ok {
			// The server closed the channel.
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		data := []byte(strconv.Itoa(message))

		err := c.conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			break
		}
	}
}

func serveWs(h *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:  h,
		conn: conn,
		send: make(chan int, 256),
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
