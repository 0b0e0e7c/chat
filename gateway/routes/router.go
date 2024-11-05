package routes

import (
	"chatting/gateway/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoutes(r *gin.Engine) {
	// User routes
	UserRoutes(r)

	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	// User routes with JWT middleware
	jwtMiddleware := middleware.NewJWTMiddleware()
	UserRoutesJWT(r, jwtMiddleware)

	// Friend routes
	FriendRoutes(r, jwtMiddleware)

	// Message routes
	MessageRoutes(r, jwtMiddleware)

	// Group message routes
	GroupMessageRoutes(r, jwtMiddleware)

	// Swagger routes

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
