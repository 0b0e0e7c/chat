package routes

import (
	"github.com/gin-gonic/gin"

	"chatting/gateway/handler"
)

func UserRoutes(r *gin.Engine) {
	userGroup := r.Group("/api/user")

	userHandler := handler.GlobalLogic
	{
		userGroup.POST("/register", userHandler.Register)
		userGroup.POST("/login", userHandler.Login)
	}
}

func UserRoutesJWT(r *gin.Engine, jwtHandler gin.HandlerFunc) {
	userGroup := r.Group("/api/user")

	routeHandler := handler.GlobalLogic
	userGroup.POST("/logout", jwtHandler, routeHandler.Logout)

	{
		profileGroup := userGroup.Group("/profile")
		profileGroup.Use(jwtHandler)
		{
			profileGroup.GET("/", routeHandler.GetProfile)
			profileGroup.PUT("/update", routeHandler.UpdateProfile)

			avatarGroup := profileGroup.Group("/avatar")
			{
				avatarGroup.GET("/", routeHandler.GetAvatar)
				avatarGroup.PUT("/upload", routeHandler.UploadAvatar)
			}
		}
	}
}
