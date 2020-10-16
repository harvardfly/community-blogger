package repositories

import (
	"community-blogger/internal/pkg/database"
	"community-blogger/internal/pkg/models"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

// UserRepository Repository 用户RPC数据库操作
type UserRepository interface {
	FindByToken(token string) (*models.User, error)
	FindByID(id int32) (models.User, error)
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

// FindByToken 数据库实现根据Token获取用户信息
func (r *MysqlUserRepository) FindByToken(token string) (*models.User, error) {
	var res *models.User
	if err := r.db.Where("token=?", token).First(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// FindByID 数据库实现根据ID获取用户信息
func (r *MysqlUserRepository) FindByID(id int32) (models.User, error) {
	var res models.User
	if err := r.db.Model(&models.User{}).Scopes(withID(id)).First(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}
