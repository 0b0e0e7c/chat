package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource string
	CacheRedis CacheRedis
	CertConfig KeyPair
}

type CacheRedis struct {
	Host string
	Type string
	Pass string
}

// 配置证书路径
type KeyPair struct {
	CertFile string
	KeyFile  string
}
