package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var count = 0

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "index.html")
}

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

		data := []byte(fmt.Sprint(count))
		err = conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
