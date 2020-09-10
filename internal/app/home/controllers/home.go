package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"community-blogger/internal/app/home/services"
	"community-blogger/internal/pkg/requests"
	"community-blogger/internal/pkg/utils/httputil"
)

// HomeController 定义HomeController结构
type HomeController struct {
	logger  *zap.Logger
	service services.HomeService
}

// NewHomeController 初始化HomeController
func NewHomeController(logger *zap.Logger, s services.HomeService) *HomeController {
	return &HomeController{
		logger:  logger.With(zap.String("type", "HomeController")),
		service: s,
	}
}

// List 获取home列表 分页
func (pc *HomeController) List(c *gin.Context) {
	var req requests.HomeList
	if err := c.ShouldBindQuery(&req); err != nil {
		pc.logger.Error("参数错误", zap.Error(err))
		c.JSON(http.StatusBadRequest, httputil.Error(nil, "参数校验失败"))
		return
	}
	count, list, err := pc.service.HomeList(&req)
	if err != nil {
		pc.logger.Error("获取home列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, httputil.Error(nil, "获取home列表失败"))
		return
	}
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["code"] = http.StatusOK
	result["data"] = list
	result["count"] = count
	c.JSON(http.StatusOK, result)
}

// Home 新增home信息
func (pc *HomeController) Home(c *gin.Context) {
	var req requests.Home
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.Error(err, "参数校验失败"))
		return
	}
	err := pc.service.Home(&req)
	if err != nil {
		pc.logger.Error("添加失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, httputil.Error(err, "添加失败"))
		return
	}
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["code"] = http.StatusOK
	c.JSON(http.StatusOK, result)
}
