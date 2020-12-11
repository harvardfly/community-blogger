package services

import (
	"community-blogger/internal/app/home/repositories"
	"community-blogger/internal/pkg/requests"
	"community-blogger/internal/pkg/responses"
	minio2 "community-blogger/internal/pkg/storages/minio"
	qiniu2 "community-blogger/internal/pkg/storages/qiniu"
	"community-blogger/internal/pkg/utils/constutil"

	"github.com/qiniu/api.v7/v7/storage"

	"github.com/minio/minio-go/v6"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// HomeService home模块service定义
type HomeService interface {
	HomeList(req *requests.HomeList) (int, []responses.Home, error)
	Home(req *requests.Home) error
	UploadFileMinio(uploadDir, filename string) (string, error)
	UploadFileQiniu(key, filename string) (string, error)
	QiniuFileInfo(key string) (responses.FileInfo, error)
}

// DefaultHomeService home模块service默认对象
type DefaultHomeService struct {
	logger        *zap.Logger
	v             *viper.Viper
	MinioClient   *minio.Client
	QiniuUploader *storage.FormUploader
	Repository    repositories.HomeRepository
}

// NewHomeService 初始化
func NewHomeService(
	logger *zap.Logger,
	v *viper.Viper,
	repository repositories.HomeRepository,
	minioClient *minio.Client,
	qiniuUploader *storage.FormUploader,
) HomeService {
	return &DefaultHomeService{
		logger:        logger.With(zap.String("type", "DefaultHomeService")),
		v:             v,
		Repository:    repository,
		MinioClient:   minioClient,
		QiniuUploader: qiniuUploader,
	}
}

// HomeList home信息列表
func (s *DefaultHomeService) HomeList(req *requests.HomeList) (int, []responses.Home, error) {
	if req.Page < 0 {
		req.Page = 0
	}
	if req.Limit <= 0 {
		req.Limit = constutil.DefaultPageSize
	} else if req.Limit > constutil.MaxPageSize {
		req.Limit = constutil.MaxPageSize
	}
	req.Paginator = map[string]int{"limit": req.Limit, "offset": req.Page}
	req.Orders = map[string]string{"order": req.Order, "desc": req.Description}
	return s.Repository.HomeList(req)
}

// Home 新建home信息
func (s *DefaultHomeService) Home(req *requests.Home) error {
	return s.Repository.Home(req)
}

// UploadFileMinio 上传文件到Minio
func (s *DefaultHomeService) UploadFileMinio(uploadDir, filename string) (string, error) {
	path, err := minio2.UploadFile(uploadDir, filename)
	if err != nil {
		s.logger.Error("", zap.Error(err))
		return "", err
	}
	return path, nil
}

// UploadFileQiniu 上传文件到qiniu
func (s *DefaultHomeService) UploadFileQiniu(key, filename string) (string, error) {
	path, err := qiniu2.UploadFile(key, filename)
	if err != nil {
		s.logger.Error("", zap.Error(err))
		return "", err
	}
	return path, nil
}

// GetQiniuFileInfo 获取qiniu文件信息
func (s *DefaultHomeService) QiniuFileInfo(key string) (responses.FileInfo, error) {
	fileInfo, err := qiniu2.GetFileInfo(key)
	if err != nil {
		s.logger.Error("获取文件信息失败", zap.Error(err))
		return responses.FileInfo{}, err
	}
	return responses.FileInfo{
		Hash:     fileInfo.Hash,
		Fsize:    fileInfo.Fsize,
		PutTime:  fileInfo.PutTime,
		MimeType: fileInfo.MimeType,
		Type:     fileInfo.Type,
	}, nil
}
