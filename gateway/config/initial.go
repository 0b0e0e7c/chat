package config

import (
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"

	"chatting/gateway/handler"
	"chatting/service/friendService/pb/friend"
	"chatting/service/groupMessageService/pb/groupMessage"
	"chatting/service/messageService/pb/message"
	"chatting/service/userService/pb/user"
)

func InitRpcClients(config *Config) error {
	var err error

	handler.UserRpcClient, err = initUserRpcClient(config.UserRpcEtcd)
	if err != nil {
		logx.Errorf("config.Init err: %v", err)
		return err
	}
	handler.FriendRpcClient, err = initFriendRpcClient(config.FriendRpcEtcd)
	if err != nil {
		logx.Errorf("config.Init err: %v", err)
		return err
	}
	handler.MsgRpcClient, err = initMsgRpcClient(config.MsgRpcEtcd)
	if err != nil {
		logx.Errorf("config.Init err: %v", err)
		return err
	}
	handler.GroupMsgRpcClient, err = initGroupMsgRpcClient(config.GroupMsgRpcEtcd)
	if err != nil {
		logx.Errorf("config.Init err: %v", err)
		return err
	}

	handler.GlobalLogic = handler.NewLogic()
	return nil
}

func initUserRpcClient(c discov.EtcdConf) (user.UserServiceClient, error) {
	userClient, err := zrpc.NewClient(
		zrpc.RpcClientConf{
			Etcd: c,
		})
	if err != nil {
		return nil, err
	}
	return user.NewUserServiceClient(userClient.Conn()), nil
}

func initFriendRpcClient(c discov.EtcdConf) (friend.FriendServiceClient, error) {
	friendClient, err := zrpc.NewClient(
		zrpc.RpcClientConf{
			Etcd: c,
		})
	if err != nil {
		return nil, err
	}
	return friend.NewFriendServiceClient(friendClient.Conn()), nil
}

func initMsgRpcClient(c discov.EtcdConf) (message.MessageServiceClient, error) {
	messageClient, err := zrpc.NewClient(
		zrpc.RpcClientConf{
			Etcd: c,
		})
	if err != nil {
		return nil, err
	}
	return message.NewMessageServiceClient(messageClient.Conn()), nil
}

func initGroupMsgRpcClient(c discov.EtcdConf) (groupMessage.GroupMessageServiceClient, error) {
	groupMsgClient, err := zrpc.NewClient(
		zrpc.RpcClientConf{
			Etcd: c,
		})
	if err != nil {
		return nil, err
	}
	return groupMessage.NewGroupMessageServiceClient(groupMsgClient.Conn()), nil
}
