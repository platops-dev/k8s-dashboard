package main

import (
	"test4/config"
	"test4/controller"
	"test4/service"

	"github.com/gin-gonic/gin"
)

func main() {
	//初始化gin对象
	r := gin.Default()
	//初始化k8s client
	service.K8s.Init()
	//初始化路由规则
	controller.Router.InitApiRouter(r)
	//gin程序启动
	r.Run(config.ListenAddr)
}