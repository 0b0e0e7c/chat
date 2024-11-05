package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	DataSource string
	CacheRedis CacheRedis
}

type CacheRedis struct {
	Host string
	Type string
	Pass string
}
