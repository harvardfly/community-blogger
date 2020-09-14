package controllers

import (
	"community-blogger/internal/pkg/transports/http"
	"community-blogger/internal/pkg/transports/http/middlewares/auth"
	"community-blogger/internal/pkg/transports/http/middlewares/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// CreateInitControllersFn  article controllers router
func CreateInitControllersFn(pc *ArticleController) http.InitControllers {
	return func(r *gin.Engine) {
		r.Use(auth.ValidAccessToken) // JWT Token认证
		r.Use(auth.Next())           // 允许跨域
		e := r.Group("api/v1")
		e.POST("/article", middleware.BucketLimit, pc.Article)   // 发表文章
		e.GET("/article", pc.ArticleInfo)                    // 获取文章详情
		e.GET("/article/topN", pc.ArticleTopN)               // 获取TOPN的文章
		e.PUT("/article", pc.ArticleEdit)                    // 更新文章  全量更新
		e.PATCH("/article/category", pc.ArticleCategoryEdit) // 修改文章类别  部分更新
		e.DELETE("article", pc.ArticleDel)                   // 删除文章
	}
}

// ProviderSet article controllers wire
var ProviderSet = wire.NewSet(NewArticleController, CreateInitControllersFn)
