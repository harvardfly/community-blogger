package oss

import (
	"community-blogger/internal/pkg/log"
	"community-blogger/internal/pkg/responses"
	"community-blogger/internal/pkg/utils/fileutil"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// ClientType 定义oss client 结构体
type ClientType struct {
	OssClient *oss.Client
	Settings  *Options
}

// Client  oss连接类型
var Client ClientType

// Options oss option
type Options struct {
	AccessKeyID     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	Bucket          string `yaml:"bucket"`
	Endpoint        string `yaml:"endpoint"`
	Domain          string `yaml:"domain"`
}

// NewOptions 加载oss配置
func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("oss", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal oss option error")
	}
	Client.Settings = o
	logger.Info("load oss options success", zap.Any("oss options", o))
	return o, err
}

// New 初始化oss客户端
func New(o *Options) (*oss.Client, error) {
	client, err := oss.New(o.Endpoint, o.AccessKeyID, o.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	// 创建存储空间
	err = client.CreateBucket(o.Bucket)
	if err != nil {
		return nil, err
	}
	Client.OssClient = client
	return client, nil
}

// UploadFile 上传文件 并返回文件信息
func UploadFile(objectName, localFilePath string) (responses.FileInfo, error) {
	/*
		objectName: 可定义为上传的文件名
		localFilePath： 文件所在的路径 包含文件名
	*/
	// 获取存储空间
	bucket, err := Client.OssClient.Bucket(Client.Settings.Bucket)
	if err != nil {
		log.Client.Logger.Error(
			"获取存储空间失败",
			zap.Error(err),
		)
		return responses.FileInfo{}, err
	}

	// 打开文件
	f, err := os.Open(localFilePath)
	if err != nil {
		log.Client.Logger.Error(
			"打开文件失败",
			zap.Error(err),
		)
		return responses.FileInfo{}, err
	}
	defer f.Close()

	// 获取文件的contentType
	contentType, err := fileutil.GetFileContentType(f)
	if err != nil {
		log.Client.Logger.Error(
			"获取文件contentType失败",
			zap.Error(err),
		)
		return responses.FileInfo{}, err
	}
	// 设置文件contentType 防止访问图片地址自动下载
	option := oss.ContentType(contentType)

	// 指定存储类型为标准存储，缺省也为标准存储。
	storageType := oss.ObjectStorageClass(oss.StorageStandard)

	// 指定访问权限为公共读，缺省为继承bucket的权限。
	objectAcl := oss.ObjectACL(oss.ACLPublicRead)

	// 上传文件
	err = bucket.PutObjectFromFile(
		objectName,
		localFilePath,
		option,
		storageType,
		objectAcl,
	)
	if err != nil {
		log.Client.Logger.Error(
			"上传失败",
			zap.Error(err),
		)
		return responses.FileInfo{}, err
	}

	// 获取文件元信息。
	props, err := bucket.GetObjectDetailedMeta(objectName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	contentLength := props.Get("Content-Length")
	fsize, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return responses.FileInfo{}, err
	}

	return responses.FileInfo{
		Hash:     props.Get("X-Oss-Hash-Crc64ecma"),
		Fsize:    fileutil.FormatFileSize(fsize),
		PutTime:  time.Now().Unix(),
		MimeType: props.Get("Content-Type"),
		URI:      fmt.Sprintf("%s/%s", Client.Settings.Domain, objectName),
	}, nil
}

// DownloadFile 下载文件
func DownloadFile(objectName, downloadedFileName string) (string, error) {
	// 获取存储空间
	bucket, err := Client.OssClient.Bucket(Client.Settings.Bucket)
	if err != nil {
		log.Client.Logger.Error(
			"获取存储空间失败",
			zap.Error(err),
		)
		return "", err
	}

	// 下载文件到服务器
	err = bucket.GetObjectToFile(objectName, objectName)
	if err != nil {
		log.Client.Logger.Error(
			"下载文件失败",
			zap.Error(err),
		)
		return "", err
	}
	return "ok", nil
}

// FileList 获取文件列表信息
func FileList() (string, error) {
	// 获取存储空间
	bucket, err := Client.OssClient.Bucket(Client.Settings.Bucket)
	if err != nil {
		log.Client.Logger.Error(
			"获取存储空间失败",
			zap.Error(err),
		)
		return "", err
	}

	// 默认列举100个文件
	lsRes, err := bucket.ListObjects()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 打印结果。
	for _, object := range lsRes.Objects {
		fmt.Println(fileutil.FormatFileSize(object.Size))
		fmt.Println("Object:", object.XMLName)
		fmt.Println(object.Type)
	}
	return "ok", nil
}

// DeleteFile 删除文件
func DeleteFile(objectName string) (string, error) {
	// 获取存储空间
	bucket, err := Client.OssClient.Bucket(Client.Settings.Bucket)
	if err != nil {
		log.Client.Logger.Error(
			"获取存储空间失败",
			zap.Error(err),
		)
		return "", err
	}

	// 删除单个文件。objectName表示删除OSS文件时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	// 如需删除文件夹，请将objectName设置为对应的文件夹名称。如果文件夹非空，则需要将文件夹下的所有object删除后才能删除该文件夹
	err = bucket.DeleteObject(objectName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	return "ok", nil
}

// DeleteFiles 删除多个文件
func DeleteFiles(objectNames []string) (string, error) {
	// 获取存储空间
	bucket, err := Client.OssClient.Bucket(Client.Settings.Bucket)
	if err != nil {
		log.Client.Logger.Error(
			"获取存储空间失败",
			zap.Error(err),
		)
		return "", err
	}

	// 返回删除成功的文件
	delRes, err := bucket.DeleteObjects(objectNames)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	fmt.Println("Deleted Objects:", delRes.DeletedObjects)

	// 不返回删除的结果
	//_, err = bucket.DeleteObjects(
	//	objectNames,
	//	oss.DeleteObjectsQuiet(true),
	//)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	os.Exit(-1)
	//}
	return "ok", nil
}

// ProviderSet inject
var ProviderSet = wire.NewSet(New, NewOptions)
