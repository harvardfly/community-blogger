package controllers

import (
	"community-blogger/internal/pkg/transports/http"
	"community-blogger/internal/pkg/transports/http/middlewares/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// CreateInitControllersFn  user controllers router
func CreateInitControllersFn(pc *UserController) http.InitControllers {
	return func(r *gin.Engine) {
		r.Use(auth.Next()) // 允许跨域
		e := r.Group("api/v1")
		// JWT Token认证
		e.GET("/user", auth.ValidAccessToken, pc.UserInfoByID)
		e.POST("/register", pc.Register) // 注册
		e.POST("/login", pc.Login)       // 登录
	}
}

// ProviderSet user controllers wire
var ProviderSet = wire.NewSet(NewUserController, CreateInitControllersFn)
