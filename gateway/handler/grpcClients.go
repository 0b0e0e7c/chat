package handler

import (
	"chatting/component/bucket"
	"chatting/service/friendService/pb/friend"
	"chatting/service/groupMessageService/pb/groupMessage"
	"chatting/service/messageService/pb/message"
	"chatting/service/userService/pb/user"
)

var (
	UserRpcClient     user.UserServiceClient
	FriendRpcClient   friend.FriendServiceClient
	MsgRpcClient      message.MessageServiceClient
	GroupMsgRpcClient groupMessage.GroupMessageServiceClient

	GlobalLogic *Logic
)

type Logic struct {
	userRpc     user.UserServiceClient
	friendRpc   friend.FriendServiceClient
	msgRpc      message.MessageServiceClient
	groupMsgRpc groupMessage.GroupMessageServiceClient

	bucketClient bucket.Minio
}

func NewLogic() *Logic {
	return &Logic{
		userRpc:      UserRpcClient,
		friendRpc:    FriendRpcClient,
		msgRpc:       MsgRpcClient,
		groupMsgRpc:  GroupMsgRpcClient,
		bucketClient: bucket.NewMinioClient(),
	}
}

func (l *Logic) UserRpc() user.UserServiceClient {
	return UserRpcClient
}

func (l *Logic) FriendRpc() friend.FriendServiceClient {
	return FriendRpcClient
}

func (l *Logic) MsgRpc() message.MessageServiceClient {
	return MsgRpcClient
}

func (l *Logic) GroupMsgRpc() groupMessage.GroupMessageServiceClient {
	return GroupMsgRpcClient
}
