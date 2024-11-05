package routes

import (
	"github.com/gin-gonic/gin"

	"chatting/gateway/handler"
)

func FriendRoutes(r *gin.Engine, jwtHandler gin.HandlerFunc) {
	friendGroup := r.Group("/api/friend")

	friendHandler := handler.GlobalLogic
	friendGroup.Use(jwtHandler)
	{
		friendGroup.POST("/pending", friendHandler.PendingFriends)
		friendGroup.POST("/add", friendHandler.Confirm)
		friendGroup.POST("/delete", friendHandler.DeleteFriend)

		friendGroup.GET("/pending", friendHandler.GetPendingFriends)
		friendGroup.GET("/friends", friendHandler.GetFriends)
	}

	friendManageGroup := friendGroup.Group("/manage")
	{
		friendManageGroup.POST("/createGroup", friendHandler.CreateFriendGroup)
		friendManageGroup.POST("/addToGroup", friendHandler.AddToGroup)
		friendManageGroup.POST("/deleteFromGroup", friendHandler.DeleteFromGroup)

		friendManageGroup.GET("/friendGroup", friendHandler.GetFriendGroup)
	}
}
