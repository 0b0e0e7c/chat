package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"

	"chatting/gateway/config"
	_ "chatting/gateway/docs"
	"chatting/gateway/routes"
)

//	@title			Jackey Chatting API
//	@version		1
//	@description	This is a chatting API

//	@contact.name	API Support
//	@contact.url	http://www.example.com/support
//	@contact.email	support@example.com

//	@BasePath	/api/
//	@host		localhost:8888

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	err := config.InitRpcClients(&c)
	if err != nil {
		logx.Errorf("init rpc clients err: %v", err)
		// maybe exit here
		os.Exit(1)
	}

	r := gin.Default()

	routes.InitRoutes(r)

	if err := r.Run(c.Listen); err != nil {
		fmt.Println(err)
		return
	}
}
