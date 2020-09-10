package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"community-blogger/internal/app/user/services"
	"community-blogger/internal/pkg/baseresponse"
	"community-blogger/internal/pkg/requests"
	"community-blogger/internal/pkg/utils/httputil"
)

// UserController 定义user模块Controller
type UserController struct {
	logger  *zap.Logger
	service services.UserService
}

// NewUserController 初始化用户模块Controller
func NewUserController(logger *zap.Logger, s services.UserService) *UserController {
	return &UserController{
		logger:  logger.With(zap.String("type", "UserController")),
		service: s,
	}
}

// UserInfoByID 调用rpc服务获取用户详情
func (pc *UserController) UserInfoByID(c *gin.Context) {
	var req requests.User
	if err := c.ShouldBindQuery(&req); err != nil {
		pc.logger.Error("参数错误", zap.Error(err))
		c.JSON(http.StatusBadRequest, httputil.Error(nil, "参数校验失败"))
		return
	}
	res, err := pc.service.FindByID(c, &req)
	if err != nil {
		pc.logger.Error("查找用户结果失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, httputil.Error(nil, "查找用户结果失败"))
		return
	}
	result := make(map[string]interface{})
	result["data"] = res
	c.JSON(http.StatusOK, result)
}

// Register 新增user信息
func (pc *UserController) Register(c *gin.Context) {
	req := new(requests.RegisterRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		baseresponse.ParamError(c, err)
		return
	}
	res, err := pc.service.Register(req)
	baseresponse.HTTPResponse(c, res, err)
	return
}

// Login 用户登录
func (pc *UserController) Login(c *gin.Context) {
	req := new(requests.LoginRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		// 通用参数验证方法
		baseresponse.ParamError(c, err)
		return
	}
	res, err := pc.service.Login(req)
	baseresponse.HTTPResponse(c, res, err)
	return
}
