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

type update struct {
	coord
	Color int `json:"color"`
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

func togglePixel(c *gin.Context) {
	var co coord

	if err := c.Bind(&co); err != nil {
		return
	}

	val, err := board.togglePixel(co.Row, co.Col)
	if err != nil {
		log.Println("togglePixel:", err)
		c.String(http.StatusServiceUnavailable, "db error")
		return
	}

	c.JSON(http.StatusOK, val)

	// inform the websocket server
	upd := update{
		coord: co,
		Color: val,
	}
	data, err := json.Marshal(upd)
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
