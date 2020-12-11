package qiniu

import (
	"context"
	"fmt"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// ClientType 定义minio client 结构体
type ClientType struct {
	FormUploader  *storage.FormUploader
	UpToken       string
	BucketManager *storage.BucketManager
	Settings      *Options
}

// Client  qiniu连接类型
var Client ClientType

// Options qiniu option
type Options struct {
	AccessKey     string
	SecretKey     string
	Bucket        string
	Zone          string
	UseHTTPS      bool
	UseCdnDomains bool
	Domain        string
}

// FileInfo 文件信息
type FileInfo struct {
	Hash     string `json:"hash"`
	Fsize    int64  `json:"fsize"`
	PutTime  int64  `json:"putTime"`
	MimeType string `json:"mimeType"`
	Type     int    `json:"type"`
}

// NewOptions 加载qiniu配置
func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("qiniu", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal qiniu option error")
	}

	logger.Info("load qiniu options success", zap.Any("qiniu options", o))
	Client.Settings = o
	return o, err
}

// getZone 根据配置文件获取对应的zone
func getZone(o *Options) *storage.Region {
	zone := &storage.Region{}
	// 空间对应的机房
	switch o.Zone {
	case "ZoneHuanan":
		zone = &storage.ZoneHuanan
	case "ZoneHuabei":
		zone = &storage.ZoneHuabei
	case "ZoneBeimei":
		zone = &storage.ZoneBeimei
	case "ZoneXinjiapo":
		zone = &storage.ZoneXinjiapo
	default:
		zone = &storage.ZoneHuadong
	}
	return zone
}

// New 初始化qiniu客户端
func New(o *Options) (*storage.FormUploader, error) {
	fmt.Println(o.Bucket)
	putPolicy := storage.PutPolicy{
		Scope: o.Bucket,
	}
	mac := qbox.NewMac(o.AccessKey, o.SecretKey)
	fmt.Println(o.AccessKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	//空间对应机房
	cfg.Zone = getZone(o)
	//不启用HTTPS域名
	cfg.UseHTTPS = o.UseHTTPS
	//不使用CND加速
	cfg.UseCdnDomains = o.UseCdnDomains
	//构建上传表单对象
	formUploader := storage.NewFormUploader(&cfg)
	bucketManager := storage.NewBucketManager(mac, &cfg)
	Client.FormUploader = formUploader
	Client.UpToken = upToken
	Client.BucketManager = bucketManager
	return formUploader, nil
}

// UploadFile 上传
func UploadFile(key, filePath string) (string, error) {
	formUploader := Client.FormUploader
	upToken := Client.UpToken
	ret := storage.PutRet{}
	// 可选
	putExtra := storage.PutExtra{}
	err := formUploader.PutFile(
		context.Background(),
		&ret,
		upToken,
		key,
		filePath,
		&putExtra,
	)
	if err != nil {
		return "", err
	}

	//返回地址
	return fmt.Sprintf("%s/%s", Client.Settings.Domain, ret.Key), nil
}

// GetFileInfo 获取文件信息
func GetFileInfo(key string) (*storage.FileInfo, error) {
	fileInfo, sErr := Client.BucketManager.Stat(Client.Settings.Bucket, key)
	if sErr != nil {
		return &storage.FileInfo{}, sErr
	}
	return &fileInfo, nil
}

// ProviderSet inject
var ProviderSet = wire.NewSet(New, NewOptions)
