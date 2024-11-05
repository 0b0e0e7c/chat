package routes

import (
	"github.com/gin-gonic/gin"

	"chatting/gateway/handler"
)

func GroupMessageRoutes(r *gin.Engine, jwtHandler gin.HandlerFunc) {
	groupMessageHandler := handler.NewGroupMessageHandler()

	groupMessageGroup := r.Group("/api/groupMessage")
	groupMessageGroup.Use(jwtHandler, gin.Logger())
	{
		groupMessageGroup.POST("/send", groupMessageHandler.SendMsg)
		groupMessageGroup.GET("/", groupMessageHandler.GetMsg)
		groupMessageGroup.GET("/ws", groupMessageHandler.WebSocketHandler)

		groupMessageGroup.POST("/create", groupMessageHandler.CreateMsgGroup)
		groupMessageGroup.POST("/addGroupMember", groupMessageHandler.AddGroupMember)
	}
}
