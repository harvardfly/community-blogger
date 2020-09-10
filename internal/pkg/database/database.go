package database

import (
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	// mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
	"community-blogger/internal/pkg/models"
)

// Options database option
type Options struct {
	Mysql struct {
		URL   string `yaml:"url"`
		Debug bool
	}
}

// NewOptions new database option
func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("database", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal database option error")
	}

	logger.Info("load database options success", zap.Any("database options", o))
	return o, err
}

// Database 定义数据库struct
type Database struct {
	Mysql *gorm.DB
}

// DBClient  mysql连接类型
var DBClient Database

// New new database
func New(o *Options) (*Database, error) {
	var d = new(Database)
	if o.Mysql.URL == "" {
		return nil, errors.New("缺少mysql配置")
	}
	if o.Mysql.URL != "" {
		mysql, err := mysql(o)
		if err != nil {
			return nil, err
		}
		d.Mysql = mysql
	}
	DBClient.Mysql = d.Mysql
	return d, nil
}

// mysql 定义mysql连接信息
func mysql(o *Options) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", o.Mysql.URL)
	if err != nil {
		return nil, errors.Wrap(err, "gorm open mysql connection error")
	}
	err = db.DB().Ping()
	if err != nil {
		return nil, errors.Wrap(err, "mysql ping fail")
	}
	if o.Mysql.Debug {
		db = db.Debug()
	}
	db.DB().SetConnMaxLifetime(time.Minute * 10)

	// 自动迁移模式
	db.AutoMigrate(&models.Home{}, &models.Article{}, &models.Category{}, &models.User{})

	return db, nil
}

// ProviderSet dependency injection
var ProviderSet = wire.NewSet(New, NewOptions)
