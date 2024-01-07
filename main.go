package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var count = 0

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		switch string(message) {
		case "increment":
			count++
		case "getCount":
		}

		data := []byte(strconv.Itoa(count))
		err = conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", serveWs)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
