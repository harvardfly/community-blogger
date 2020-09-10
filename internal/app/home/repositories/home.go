package repositories

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"time"
	"community-blogger/internal/pkg/database"
	"community-blogger/internal/pkg/models"
	"community-blogger/internal/pkg/requests"
	"community-blogger/internal/pkg/responses"
	"community-blogger/internal/pkg/utils/dbutil"
)

// HomeRepository home模块Repository
type HomeRepository interface {
	HomeList(req *requests.HomeList) (int, []responses.Home, error)
	Home(req *requests.Home) error
}

// MysqlHomeRepository home模块MysqlHomeRepository
type MysqlHomeRepository struct {
	logger *zap.Logger
	db     *gorm.DB
}

// NewMysqlHomeRepository 初始化MysqlHomeRepository
func NewMysqlHomeRepository(logger *zap.Logger, db *database.Database) HomeRepository {
	return &MysqlHomeRepository{
		logger: logger.With(zap.String("type", "MysqlHomeRepository")),
		db:     db.Mysql,
	}
}

// HomeList home信息列表
func (r *MysqlHomeRepository) HomeList(req *requests.HomeList) (int, []responses.Home, error) {
	var (
		count int
		value []responses.Home
	)
	r.db.Model(&models.Home{}).Where("description like ?", "%"+req.Description+"%").Count(&count)
	sqlStr := dbutil.PageOptimize(req.Paginator, req.Orders, models.Home{}.TableName(), " where description like "+"'%"+req.Description+"%'")
	err := r.db.Model(&models.Home{}).Raw(sqlStr).Find(&value).Error
	return count, value, err
}

// Home 创建home信息
func (r *MysqlHomeRepository) Home(req *requests.Home) error {
	var (
		home models.Home
	)
	home.URL = req.URL
	home.Img = req.Img
	home.Title = req.Title
	home.Description = req.Description
	home.CreatedAt = time.Now()
	home.UpdatedAt = time.Now()
	err := r.db.Model(models.Home{}).Create(&home).Error
	if err != nil {
		return err
	}
	return nil
}
