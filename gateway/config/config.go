package config

import "github.com/zeromicro/go-zero/core/discov"

type Config struct {
	Listen          string          `yaml:"listen" json:"listen"`
	UserRpcEtcd     discov.EtcdConf `yaml:"userRpcEtcd" json:"userRpcEtcd"`
	FriendRpcEtcd   discov.EtcdConf `yaml:"friendRpcEtcd" json:"friendRpcEtcd"`
	MsgRpcEtcd      discov.EtcdConf `yaml:"msgRpcEtcd" json:"msgRpcEtcd"`
	GroupMsgRpcEtcd discov.EtcdConf `yaml:"groupMsgRpcEtcd" json:"groupMsgRpcEtcd"`
}
