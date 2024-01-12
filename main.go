package main

import (
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

	err := c.Bind(&pos)
	if err != nil {
		return
	}

	i := pos.X*width + pos.Y
	// TODO: Handle OOB
	v := (grid[i/8] >> (7 - i%8)) & 1

	c.JSON(http.StatusOK, v)
}

func main() {
	router := gin.Default()
	router.GET("/pixel", getPixel)

	router.Run(":8080")
}
