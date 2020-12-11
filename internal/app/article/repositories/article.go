package repositories

import (
	"community-blogger/internal/pkg/database"
	"community-blogger/internal/pkg/kafka"
	"community-blogger/internal/pkg/models"
	"community-blogger/internal/pkg/requests"
	"community-blogger/internal/pkg/responses"
	"community-blogger/internal/pkg/utils/constutil"
	"context"
	"errors"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/Shopify/sarama"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

// ArticleRepository Repository 文章模块数据库操作
type ArticleRepository interface {
	Article(req *requests.Article, ctx context.Context) (models.Article, error)
	GetArticle(id int) models.Article
	ArticleReadCount(id int) error
	GetTOPNArticles(ids []int) (value []models.Article, err error)
	ArticleEdit(article *models.Article) error
	GetCategory(id int) models.Category
	ArticleCategoryEdit(article *models.Article) error
	ArticleDel(id int) error
}

// MysqlArticleRepository 定义article repository 初始结构体
type MysqlArticleRepository struct {
	logger   *zap.Logger
	db       *gorm.DB
	producer sarama.SyncProducer
}

// NewMysqlArticleRepository 初始化数据库操作
func NewMysqlArticleRepository(logger *zap.Logger, db *database.Database, producer sarama.SyncProducer) ArticleRepository {
	return &MysqlArticleRepository{
		logger:   logger.With(zap.String("type", "MysqlArticleRepository")),
		db:       db.Mysql,
		producer: producer,
	}
}

// withID 通过单个ID查询对应记录
func withID(id int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

// withIds 通过一组ID查询对应记录
func withIds(ids []int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id in (?)", ids)
	}
}

// GetArticle 获取文章及相应的类别信息
func (r *MysqlArticleRepository) GetArticle(id int) models.Article {
	var value models.Article
	r.db.Preload("Category").Model(&models.Article{}).Scopes(withID(id)).First(&value)
	return value
}

// GetCategory 根据类别获取类别对象信息
func (r *MysqlArticleRepository) GetCategory(id int) models.Category {
	var value models.Category
	r.db.Model(&models.Category{}).Scopes(withID(id)).First(&value)
	return value
}

// Article 发布文章
func (r *MysqlArticleRepository) Article(req *requests.Article, ctx context.Context) (models.Article, error) {
	var category models.Category
	value := models.Article{
		Title:      req.Title,
		Summary:    req.Summary,
		CategoryID: req.CategoryID,
	}
	span, _ := opentracing.StartSpanFromContext(ctx, "article.mysql")
	defer span.Finish()
	err := r.db.Model(&models.Article{}).Create(&value).Error
	if err != nil {
		return value, err
	}
	r.db.Model(&models.Category{}).Scopes(withID(value.CategoryID)).First(&category)
	esValue := requests.ArticleES{
		ID:      value.ID,
		Summary: value.Summary,
		Title:   value.Title,
		Category: responses.Category{
			ID:        value.CategoryID,
			Name:      category.Name,
			Num:       category.Num,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		},
		CreatedAt: value.CreatedAt,
		UpdatedAt: value.UpdatedAt,
	}
	//同步数据到kafka
	result := kafka.Client.SyncProduce(constutil.CreateArticle, esValue)
	r.logger.Info("sync kafka to es", zap.Any("result", result))

	span.LogFields(
		log.String("event", "insert article repositories"),
	)

	return value, err
}

// ArticleReadCount 文章浏览 验证文章ID是否存在
func (r *MysqlArticleRepository) ArticleReadCount(id int) error {
	err := r.db.Scopes(withID(id)).Error
	if err != nil {
		return err
	}
	return nil
}

// GetTOPNArticles 根据redis获取的文章ids查询数据库取文章信息
func (r *MysqlArticleRepository) GetTOPNArticles(ids []int) ([]models.Article, error) {
	var (
		value []models.Article
	)
	err := r.db.Model(&models.Article{}).Scopes(withIds(ids)).Find(&value).Error
	if err != nil {
		return value, err
	}

	return value, nil
}

// ArticleEdit 文章全量更新
func (r *MysqlArticleRepository) ArticleEdit(article *models.Article) error {
	err := r.db.Model(&models.Article{}).Scopes(withID(article.ID)).UpdateColumns(map[string]interface{}{
		"category_id": article.CategoryID,
		"summary":     article.Summary,
		"title":       article.Title,
		"updated_at":  time.Now(),
	}).Error
	if err != nil {
		r.logger.Error("更新文章失败", zap.Error(err))
		return errors.New("更新文章失败")
	}

	return nil
}

// ArticleCategoryEdit 文章类别更新
func (r *MysqlArticleRepository) ArticleCategoryEdit(article *models.Article) error {
	err := r.db.Model(&models.Article{}).Scopes(withID(article.ID)).UpdateColumns(map[string]interface{}{
		"category_id": article.CategoryID,
		"updated_at":  time.Now(),
	}).Error
	if err != nil {
		r.logger.Error("更新文章类别失败", zap.Error(err))
		return errors.New("更新文章类别失败")
	}

	return nil
}

// ArticleDel 删除文章
func (r *MysqlArticleRepository) ArticleDel(id int) error {
	err := r.db.Scopes(withID(id)).Delete(&models.Article{}).Error
	if err != nil {
		r.logger.Error("删除文章失败", zap.Error(err))
		return errors.New("删除文章失败")
	}

	return nil
}
