package server

import (
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	"github.com/gin-gonic/gin"
	"importer/api"
	"importer/app/config"
	"importer/app/log"
	negt "importer/pkg/go-sdk/service/http"
	"net/http"
)

var svc api.NegtServer

// New Server 服务层，该层封装服务级别的接口函数，
// 如http服务对外提供的url,grpc服务对外提供的proto
// New 提供服务的创建方法，在di中进行依赖注入
func New(s api.NegtServer) (srv common.Service, err error) {
	// 创建路由转发
	r := gin.Default()
	mux := http.NewServeMux()
	log.GinLog(config.GetLogConfig().GinLogPath, r)
	// 设置路由句柄
	initRoute(r)
	mux.Handle("/", r)
	// 启动服务
	srv = negt.NewServiceWithMux(fmt.Sprintf(":%d", config.GetServiceConfig().HttpPort), mux)
	svc = s // 给包变量svc赋值为初始化后的service
	return srv, err
}

// initRoute http请求路由设置
func initRoute(r *gin.Engine) {
	r.GET("/ping", pingHandler)
	r.POST("/rps", RpsHandler)
	r.POST("/fxj", FxjHandler)
}

// ping命令
func pingHandler(c *gin.Context) {
	c.JSON(200, "pong")
}

func RpsHandler(c *gin.Context) {
	date := c.PostForm("date")
	if date == "" {
		c.JSON(400, "date is empty")
		return
	}
	err := svc.ImportRps(date)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	c.String(200, "success")
}

func FxjHandler(c *gin.Context) {
	err := svc.ImportFxj()
	if err != nil {
		c.String(400, err.Error())
		return
	}
	c.String(200, "success")
}
