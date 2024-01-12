package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	width  = 16
	height = 16
)

var grid = [width * height / 8]byte{}

func getPixel(c *gin.Context) {
	var pos struct {
		X int `form:"x"`
		Y int `form:"y"`
	}

	if err := c.Bind(&pos); err != nil {
		return
	}

	i := pos.X*width + pos.Y
	// TODO: Handle OOB
	v := (grid[i/8] >> (7 - i%8)) & 1

	c.String(http.StatusOK, "%b", v)
}

func togglePixel(c *gin.Context) {
	var pos struct {
		X int `json:"x"`
		Y int `json:"y"`
	}

	if err := c.Bind(&pos); err != nil {
		return
	}

	i := pos.X*width + pos.Y
	// TODO: Handle OOB
	grid[i/8] ^= (1 << (7 - i%8))

	v := (grid[i/8] >> (7 - i%8)) & 1

	data, err := json.Marshal(struct {
		X     int  `json:"x"`
		Y     int  `json:"y"`
		Color byte `json:"color"`
	}{
		X:     pos.X,
		Y:     pos.Y,
		Color: v,
	})
	if err != nil {
		log.Println(err)
	}

	wss.broadcast <- data

	c.String(http.StatusOK, "")
}

func getBoard(c *gin.Context) {
	c.String(http.StatusOK, string(grid[:]))
}

func main() {
	router := gin.Default()
	router.GET("/pixel", getPixel)
	router.POST("/toggle", togglePixel)
	router.GET("/board", getBoard)

	router.StaticFile("/", "./index.html")
	router.StaticFile("/index.js", "./index.js")

	go wss.run()
	router.GET("/ws", serveWs)

	router.Run(":8080")
}
