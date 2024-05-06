// The API handlers for manipulating board state.

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type coord struct {
	Row int `form:"row" json:"row"`
	Col int `form:"col" json:"col"`
}

type pixel struct {
	coord
	Color int `form:"color" json:"color"`
}

func getPixel(c *gin.Context) {
	var co coord

	if err := c.Bind(&co); err != nil {
		return
	}

	val, err := board.getPixel(co.Row, co.Col)
	if err != nil {
		log.Println("getPixel:", err)
		c.String(http.StatusServiceUnavailable, "db error")
		return
	}

	c.JSON(http.StatusOK, val)
}

func setPixel(c *gin.Context) {
	var p pixel

	if err := c.Bind(&p); err != nil {
		return
	}

	err := board.setPixel(p.Row, p.Col, p.Color)
	if err != nil {
		log.Println("setPixel:", err)
		c.String(http.StatusServiceUnavailable, "db error")
		return
	}

	c.JSON(http.StatusOK, p.Color)

	// Send the update to the websocket server.
	data, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}
	wss.broadcast <- data
}

func getBoard(c *gin.Context) {
	val, err := board.getBoard()
	if err != nil {
		log.Println("getBoard:", err)
		c.String(http.StatusServiceUnavailable, "db error")
		return
	}

	c.String(http.StatusOK, val)
}
