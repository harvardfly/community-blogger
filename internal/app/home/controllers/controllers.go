package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"community-blogger/internal/pkg/transports/http"
	"community-blogger/internal/pkg/transports/http/middlewares/auth"
)

// CreateInitControllersFn 定义home模块router
func CreateInitControllersFn(ho *HomeController) http.InitControllers {
	return func(r *gin.Engine) {
		r.Use(auth.Next())
		e := r.Group("api/v1")
		e.GET("/home/list", ho.List) // home信息列表
		e.POST("/home", ho.Home)     // 添加home信息
	}
}

// ProviderSet home模块wire provider
var ProviderSet = wire.NewSet(NewHomeController, CreateInitControllersFn)
