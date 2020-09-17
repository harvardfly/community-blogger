package controllers

import (
	"community-blogger/internal/pkg/transports/http"
	"community-blogger/internal/pkg/transports/http/middlewares/auth"
	"community-blogger/internal/pkg/transports/http/middlewares/ginprometheus"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// CreateInitControllersFn 定义home模块router
func CreateInitControllersFn(ho *HomeController) http.InitControllers {
	return func(r *gin.Engine) {
		r.Use(auth.Next())
		gp := ginprometheus.New(r)
		r.Use(gp.Middleware())
		// metrics采样
		r.GET("/metrics", gin.WrapH(promhttp.Handler()))
		e := r.Group("api/v1")
		e.GET("/home/list", ho.List) // home信息列表
		e.POST("/home", ho.Home)     // 添加home信息
	}
}

// ProviderSet home模块wire provider
var ProviderSet = wire.NewSet(NewHomeController, CreateInitControllersFn)
