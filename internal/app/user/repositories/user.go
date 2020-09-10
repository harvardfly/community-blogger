package repositories

import (
	"crypto/md5"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"time"
	"community-blogger/internal/pkg/database"
	"community-blogger/internal/pkg/models"
	"community-blogger/internal/pkg/requests"
	"community-blogger/internal/pkg/responses"
)

// UserRepository Repository 用户模块数据库操作
type UserRepository interface {
	Register(req *requests.RegisterRequest) (*responses.RegisterResponse, error)
	FindByUserName(userName string) (*models.User, error)
}

// MysqlUserRepository 定义user repository 初始结构体
type MysqlUserRepository struct {
	logger *zap.Logger
	db     *gorm.DB
}

// NewMysqlUserRepository 初始化数据库操作
func NewMysqlUserRepository(logger *zap.Logger, db *database.Database) UserRepository {
	return &MysqlUserRepository{
		logger: logger.With(zap.String("type", "MysqlUserRepository")),
		db:     db.Mysql,
	}
}

// withID 通过单个ID查询对应记录
func withID(typeid int32) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if typeid > 0 {
			db = db.Where("id = ?", typeid)
		}
		return db
	}
}

// FindByUserName 根据用户名查找对应用户
func (r *MysqlUserRepository) FindByUserName(userName string) (*models.User, error) {
	user := new(models.User)
	err := r.db.Model(&models.User{}).Where("username=?", userName).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Register 注册
func (r *MysqlUserRepository) Register(req *requests.RegisterRequest) (*responses.RegisterResponse, error) {
	var (
		userCreate models.User
	)
	userCreate.Username = req.Username
	userCreate.Password = fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))
	userCreate.Nickname = req.Nickname
	userCreate.Mobile = req.Mobile
	userCreate.Token = uuid.NewV4().String()
	userCreate.CreatedAt = time.Now()
	userCreate.UpdatedAt = time.Now()
	res := &responses.RegisterResponse{
		Username: userCreate.Username,
		Token:    userCreate.Token,
	}
	err := r.db.Model(&models.User{}).Create(&userCreate).Error
	if err != nil {
		return res, err
	}
	return res, nil
}
