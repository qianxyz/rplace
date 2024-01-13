package main

import "github.com/gin-gonic/gin"

const (
	width  = 16
	height = 16
)

func main() {
	router := gin.Default()

	// API endpoints
	router.GET("/pixel", getPixel)
	router.POST("/toggle", togglePixel)
	router.GET("/board", getBoard)

	// Front end
	router.StaticFile("/", "./index.html")
	router.StaticFile("/index.js", "./index.js")

	// websocket
	go wss.run()
	router.GET("/ws", serveWs)

	router.Run(":8080")
}
