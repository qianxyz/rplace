package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// TODO: Locks
var count = 0
var conns = make(map[*websocket.Conn]bool)

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	conns[conn] = true

	defer func() {
		delete(conns, conn)
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		switch string(message) {
		case "increment":
			count++
			data := []byte(strconv.Itoa(count))

			for c := range conns {
				err = c.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					log.Println("write:", err)
					break
				}
			}
		case "getCount":
			data := []byte(strconv.Itoa(count))
			err = conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", serveWs)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
