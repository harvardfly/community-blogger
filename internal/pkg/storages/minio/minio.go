package minio

import (
	"fmt"
	"log"
	"mime"
	"os"
	"strings"

	"github.com/google/wire"
	"github.com/minio/minio-go/v6"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// ClientType 定义minio client 结构体
type ClientType struct {
	MinioClient *minio.Client
	Settings    *Options
}

// Client  minio连接类型
var Client ClientType

// Options minio option
type Options struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
	Location        string
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("minio", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal minio option error")
	}

	logger.Info("load minio options success", zap.Any("minio options", o))
	Client.Settings = o
	return o, err
}

func New(o *Options) (*minio.Client, error) {
	minioClient, err := minio.New(
		o.Endpoint,
		o.AccessKeyID,
		o.SecretAccessKey,
		o.UseSSL,
	)
	if err != nil {
		log.Fatalln(err)
	}
	// 创建存储桶
	exist, err := minioClient.BucketExists(o.BucketName)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	if !exist {
		err = minioClient.MakeBucket(o.BucketName, "")
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		err = minioClient.SetBucketPolicy(
			o.BucketName,
			fmt.Sprintf(
				`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`,
				o.BucketName,
			),
		)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
	}
	Client.MinioClient = minioClient
	return minioClient, nil
}

func UploadFile(uploadDir, tempfile string) (string, error) {
	objectName := tempfile
	filePath := uploadDir + tempfile
	//提取文件后缀类型
	var ext string
	if pos := strings.LastIndexByte(objectName, '.'); pos != -1 {
		ext = objectName[pos:]
		if ext == "." {
			ext = ""
		}
	}
	//返回文件扩展类型
	contentType := mime.TypeByExtension(ext)
	// Put Object
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	n, err := Client.MinioClient.PutObject(
		Client.Settings.BucketName,
		objectName,
		file,
		fileInfo.Size(),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
	var http string
	if Client.Settings.UseSSL == true {
		http = "https://" + Client.Settings.Endpoint + "/"
	} else {
		http = "http://" + Client.Settings.Endpoint + "/"
	}
	path := http + Client.Settings.BucketName + "/" + objectName
	return path, nil
}

var ProviderSet = wire.NewSet(New, NewOptions)
