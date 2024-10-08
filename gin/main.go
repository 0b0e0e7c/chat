package main

import (
	"fmt"

	"github.com/0b0e0e7c/chat/component/handler"
	"github.com/0b0e0e7c/chat/component/middleware"
	"github.com/0b0e0e7c/chat/service/friend-service/pb/friend"
	"github.com/0b0e0e7c/chat/service/message-service/pb/message"
	"github.com/0b0e0e7c/chat/service/user-service/pb/user"

	"github.com/gin-gonic/gin"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

func main() {
	userRPCClient := initUserRPCClient()
	friendRPCClient := initFriendRPCClient()
	msgRPCClient := initMsgRPCClient()

	r := gin.Default()

	r.POST("/api/user/register", func(ctx *gin.Context) {
		handler.Register(ctx, userRPCClient)
	})
	r.POST("/api/user/login", func(ctx *gin.Context) {
		handler.Login(ctx, userRPCClient)
	})

	r.POST("/api/user/ValidateJWT", func(ctx *gin.Context) {
		handler.ValidateJWT(ctx, userRPCClient)
	})

	friendGroup := r.Group("/api/friend")
	friendGroup.Use(middleware.JWTMiddleware(userRPCClient))
	{
		friendGroup.POST("/add", func(ctx *gin.Context) {
			handler.AddFriend(ctx, friendRPCClient)
		})
		friendGroup.GET("/get", func(ctx *gin.Context) {
			handler.GetFriends(ctx, friendRPCClient)
		})

	}

	messageGroup := r.Group("/api/message")
	messageGroup.Use(middleware.JWTMiddleware(userRPCClient))
	{
		messageGroup.POST("/send", func(ctx *gin.Context) {
			handler.SendMsg(ctx, msgRPCClient)
		})

		messageGroup.GET("/get", func(ctx *gin.Context) {
			handler.GetMsg(ctx, msgRPCClient)
		})
	}

	r.Run(":8888")
}

func initUserRPCClient() user.UserServiceClient {
	userClient, err := zrpc.NewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:2379"},
			Key:   "user.rpc",
		},
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return user.NewUserServiceClient(userClient.Conn())
}

func initFriendRPCClient() friend.FriendServiceClient {
	friendClient, err := zrpc.NewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:2379"},
			Key:   "friend.rpc",
		},
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return friend.NewFriendServiceClient(friendClient.Conn())
}

func initMsgRPCClient() message.MessageServiceClient {
	messageClient, err := zrpc.NewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:2379"},
			Key:   "message.rpc",
		},
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return message.NewMessageServiceClient(messageClient.Conn())
}
