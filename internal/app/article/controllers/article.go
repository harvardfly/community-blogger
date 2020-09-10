package controllers

import (
	"community-blogger/internal/app/article/services"
	"community-blogger/internal/pkg/jaeger"
	"community-blogger/internal/pkg/requests"
	"community-blogger/internal/pkg/utils/httputil"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
	"net/http"
)

// ArticleController 定义文章模块Controller
type ArticleController struct {
	logger  *zap.Logger
	service services.ArticleService
}

// NewArticleController 初始化文章模块Controller
func NewArticleController(logger *zap.Logger, s services.ArticleService) *ArticleController {
	return &ArticleController{
		logger:  logger.With(zap.String("type", "ArticleController")),
		service: s,
	}
}

// Article 发布文章
func (pc *ArticleController) Article(c *gin.Context) {
	var req requests.Article
	if err := c.ShouldBind(&req); err != nil {
		pc.logger.Error("参数错误", zap.Error(err))
		c.JSON(http.StatusBadRequest, httputil.Error(nil, "参数校验失败"))
		return
	}
	tracer := jaeger.Client.Tracer
	opentracing.SetGlobalTracer(tracer)
	span := tracer.StartSpan("Article")
	defer span.Finish()
	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)
	span.SetTag("http.url", c.Request.URL.Path)
	article, err := pc.service.Article(ctx, &req)
	if err != nil {
		pc.logger.Error("发表文章失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, httputil.Error(nil, "发表主题失败"))
		return
	}
	span.LogFields(
		log.String("event", "info"),
		log.Int("article", article.ID),
	)
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": article})
}

// ArticleInfo 文章详情 点击文章，阅读数加 1 ZIncrBy
func (pc *ArticleController) ArticleInfo(c *gin.Context) {
	var req requests.ArticleInfo
	if err := c.ShouldBindQuery(&req); err != nil {
		pc.logger.Error("参数错误", zap.Error(err))
		c.JSON(http.StatusBadRequest, httputil.Error(nil, "参数校验失败"))
		return
	}
	res, err := pc.service.GetArticle(&req)
	if err != nil {
		pc.logger.Error("获取文章详情失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, httputil.Error(nil, "获取文章详情失败"))
		return
	}

	if res.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "error"})
		return
	}

	// redis incr文章浏览次数 放到协程执行
	go func() {
		err = pc.service.ArticleReadCount(&req)
		if err != nil {
			pc.logger.Error("文章浏览", zap.Error(err))
			c.JSON(http.StatusInternalServerError, httputil.Error(nil, "文章浏览失败"))
			return
		}
	}()

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": res})
}

// ArticleTopN 按照阅读数排行返回前n篇文章的id和title
func (pc *ArticleController) ArticleTopN(c *gin.Context) {
	var req requests.ArticleTop
	if err := c.ShouldBindQuery(&req); err != nil {
		pc.logger.Error("参数错误", zap.Error(err))
		c.JSON(http.StatusBadRequest, httputil.Error(nil, "参数校验失败"))
		return
	}
	// 调用业务逻辑层 获取返回数据结果
	res, err := pc.service.GetArticleReadCountTopN(&req)
	if err != nil {
		pc.logger.Error("获取文章TOPN失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, httputil.Error(nil, "获取文章TOPN失败"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": res})
	return
}

// ArticleEdit 修改文章 全量更新
func (pc *ArticleController) ArticleEdit(c *gin.Context) {
	var req requests.ArticleEdit
	if err := c.ShouldBind(&req); err != nil {
		pc.logger.Error("参数错误", zap.Error(err))
		c.JSON(http.StatusBadRequest, httputil.Error(nil, "参数校验失败"))
		return
	}
	err := pc.service.ArticleEdit(&req)
	if err != nil {
		pc.logger.Error("更新文章失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, httputil.Error(nil, "更新文章失败"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}

// ArticleCategoryEdit 修改文章类别 部分更新
func (pc *ArticleController) ArticleCategoryEdit(c *gin.Context) {
	var req requests.ArticleCategoryEdit
	if err := c.ShouldBind(&req); err != nil {
		pc.logger.Error("参数错误", zap.Error(err))
		c.JSON(http.StatusBadRequest, httputil.Error(nil, "参数校验失败"))
		return
	}
	err := pc.service.ArticleCategoryEdit(&req)
	if err != nil {
		pc.logger.Error("更新文章失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, httputil.Error(nil, "更新文章类别失败"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}

// ArticleDel 删除文章
func (pc *ArticleController) ArticleDel(c *gin.Context) {
	var req requests.ArticleInfo
	if err := c.ShouldBind(&req); err != nil {
		pc.logger.Error("参数错误", zap.Error(err))
		c.JSON(http.StatusBadRequest, httputil.Error(nil, "参数校验失败"))
		return
	}
	err := pc.service.ArticleDel(&req)
	if err != nil {
		pc.logger.Error("删除文章失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}
