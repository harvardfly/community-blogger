package controllers

import (
	"community-blogger/internal/pkg/transports/http"
	"community-blogger/internal/pkg/transports/http/middlewares/auth"
	"community-blogger/internal/pkg/transports/http/middlewares/csrf"
	"community-blogger/internal/pkg/transports/http/middlewares/ginprometheus"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

		// 添加csrf
		store := cookie.NewStore([]byte("secret"))
		r.Use(sessions.Sessions("mysession", store))
		r.Use(csrf.Middleware(csrf.Options{
			Secret: "secret123!@#",
			ErrorFunc: func(c *gin.Context) {
				c.String(400, "CSRF token mismatch")
				c.Abort()
			},
		}))

		// 获取csrf token
		r.GET("/protected", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"csrf_token": csrf.GetToken(c),
			})
		})

		// metrics采样
		r.GET("/metrics", gin.WrapH(promhttp.Handler()))
		e := r.Group("api/v1")
		e.GET("/home/list", ho.List)     // home信息列表
		e.POST("/home", ho.Home)         // 添加home信息
		e.POST("/upload", ho.UploadFile) // 上传文件
		e.GET("/file/info", ho.FileInfo) // 获取文件信息
	}
}

// ProviderSet home模块wire provider
var ProviderSet = wire.NewSet(NewHomeController, CreateInitControllersFn)
