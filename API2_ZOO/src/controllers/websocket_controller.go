package controllers

import (
	"api2/websocket"
	"github.com/gin-gonic/gin"
)

func WebSocketHandler(c *gin.Context) {
	websocket.HandleConnections(c)
}
