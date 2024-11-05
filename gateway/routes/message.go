package routes

import (
	"github.com/gin-gonic/gin"

	"chatting/gateway/handler"
)

func MessageRoutes(r *gin.Engine, jwtHandler gin.HandlerFunc) {
	messageHandler := handler.NewMessageHandler()

	messageGroup := r.Group("/api/message")
	messageGroup.Use(jwtHandler)
	{
		messageGroup.POST("/send", messageHandler.SendMsg)
		messageGroup.GET("/", messageHandler.GetMsg)
		messageGroup.GET("/ws", messageHandler.WebSocketHandler)
	}
}
