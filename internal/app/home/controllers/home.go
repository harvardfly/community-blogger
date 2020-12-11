package controllers

import (
	"community-blogger/internal/app/home/services"
	"community-blogger/internal/pkg/requests"
	"community-blogger/internal/pkg/utils/httputil"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

// UploadFile 上传文件 返回文件地址
func (pc *HomeController) UploadFile(c *gin.Context) {
	fileImg, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.Error(err, "参数错误"))
		return
	}
	//重命名文件的名称
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	ti := tm.Format("2006010203040501")
	//提取文件后缀类型
	var ext string
	if pos := strings.LastIndexByte(header.Filename, '.'); pos != -1 {
		ext = header.Filename[pos:]
		if ext == "." {
			ext = ""
		}
	}
	filename := ti + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ext
	//创建文件
	uploadDir := "static/uploadfile/home/"
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.Error(err, "创建文件夹失败"))
	}
	out, err := os.Create(uploadDir + filename)

	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, fileImg)

	if err != nil {
		pc.logger.Error("上传失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, httputil.Error(err, "上传失败"))
		return
	}
	// 上传文件到minio
	//path, err := pc.service.UploadFileMinio(uploadDir, filename)

	// 上传文件到qiniu
	path, err := pc.service.UploadFileQiniu(filename, uploadDir+filename)

	if err != nil {
		pc.logger.Error("获取上传文件失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, httputil.Error(err, "获取上传文件失败"))
		return
	}
	os.Remove(uploadDir + filename)
	//定义返回的数据结构
	list := make(map[string]string)
	list["resource_url"] = path
	result := make(map[string]interface{})
	result["code"] = http.StatusOK
	result["data"] = list
	c.JSON(http.StatusOK, result)
}

// FileInfo 获取文件信息
func (pc *HomeController) FileInfo(c *gin.Context) {
	var req requests.FileInfo
	if err := c.ShouldBindQuery(&req); err != nil {
		pc.logger.Error("参数错误", zap.Error(err))
		c.JSON(http.StatusBadRequest, httputil.Error(nil, "参数校验失败"))
		return
	}
	fileInfo, err := pc.service.QiniuFileInfo(req.FileName)
	if err != nil {
		pc.logger.Error("获取文件信息", zap.Error(err))
		c.JSON(http.StatusInternalServerError, httputil.Error(nil, "获取文件信息失败"))
		return
	}
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["code"] = http.StatusOK
	result["data"] = fileInfo
	c.JSON(http.StatusOK, result)
}
