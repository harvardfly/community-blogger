package services

import (
	"community-blogger/internal/app/home/repositories"
	"community-blogger/internal/pkg/requests"
	"community-blogger/internal/pkg/responses"
	minio2 "community-blogger/internal/pkg/storages/minio"
	oss2 "community-blogger/internal/pkg/storages/oss"
	qiniu2 "community-blogger/internal/pkg/storages/qiniu"
	"community-blogger/internal/pkg/utils/constutil"
	"community-blogger/internal/pkg/utils/fileutil"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

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
	UploadFileOss(uploadDir, filename string) (responses.FileInfo, error)
	DownloadFileOss(uploadDir, filename string) (string, error)
	FileListOss() (string, error)
	DeleteFileOss(filename string) (string, error)
	DeleteFilesOss(objectNames []string) (string, error)
}

// DefaultHomeService home模块service默认对象
type DefaultHomeService struct {
	logger        *zap.Logger
	v             *viper.Viper
	MinioClient   *minio.Client
	QiniuUploader *storage.FormUploader
	OssClient     *oss.Client
	Repository    repositories.HomeRepository
}

// NewHomeService 初始化
func NewHomeService(
	logger *zap.Logger,
	v *viper.Viper,
	repository repositories.HomeRepository,
	minioClient *minio.Client,
	qiniuUploader *storage.FormUploader,
	ossClient *oss.Client,
) HomeService {
	return &DefaultHomeService{
		logger:        logger.With(zap.String("type", "DefaultHomeService")),
		v:             v,
		Repository:    repository,
		MinioClient:   minioClient,
		QiniuUploader: qiniuUploader,
		OssClient:     ossClient,
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
		Fsize:    fileutil.FormatFileSize(fileInfo.Fsize),
		PutTime:  fileInfo.PutTime,
		MimeType: fileInfo.MimeType,
	}, nil
}

// UploadFileOss 上传文件到oss
func (s *DefaultHomeService) UploadFileOss(objectName, localFilePath string) (responses.FileInfo, error) {
	res, err := oss2.UploadFile(objectName, localFilePath)
	if err != nil {
		s.logger.Error("", zap.Error(err))
		return responses.FileInfo{}, err
	}
	return res, nil
}

// DownloadFileOss oss下载文件
func (s *DefaultHomeService) DownloadFileOss(objectName, localFilePath string) (string, error) {
	path, err := oss2.DownloadFile(objectName, localFilePath)
	if err != nil {
		s.logger.Error("", zap.Error(err))
		return "", err
	}
	return path, nil
}

// FileListOss oss文件信息列表
func (s *DefaultHomeService) FileListOss() (string, error) {
	path, err := oss2.FileList()
	if err != nil {
		s.logger.Error("", zap.Error(err))
		return "", err
	}
	return path, nil
}

// DeleteFileOss 删除oss文件
func (s *DefaultHomeService) DeleteFileOss(objectName string) (string, error) {
	res, err := oss2.DeleteFile(objectName)
	if err != nil {
		s.logger.Error("", zap.Error(err))
		return "", err
	}
	return res, nil
}

// DeleteFilesOss 删除多个oss文件
func (s *DefaultHomeService) DeleteFilesOss(objectNames []string) (string, error) {
	res, err := oss2.DeleteFiles(objectNames)
	if err != nil {
		s.logger.Error("", zap.Error(err))
		return "", err
	}
	return res, nil
}
