package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

const (
	width  = 16
	height = 16
)

func main() {
	router := gin.Default()

	// API endpoints
	router.GET("/pixel", getPixel)
	router.POST("/pixel", setPixel)
	router.GET("/board", getBoard)

	// Front end
	router.Use(static.Serve("/", static.LocalFile("./static", false)))

	// websocket
	go wss.run()
	router.GET("/ws", serveWs)

	router.Run(":8080")
}
