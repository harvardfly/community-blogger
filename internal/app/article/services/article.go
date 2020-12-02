package services

import (
	"community-blogger/internal/app/article/repositories"
	"community-blogger/internal/pkg/models"
	"community-blogger/internal/pkg/redis"
	"community-blogger/internal/pkg/requests"
	"community-blogger/internal/pkg/responses"
	"community-blogger/internal/pkg/utils/constutil"
	"community-blogger/internal/pkg/utils/timeutil"
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go/log"

	redisPool "github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// ArticleService 定义article service
type ArticleService interface {
	Article(ctx context.Context, req *requests.Article) (responses.Article, error)
	GetArticle(req *requests.ArticleInfo) (responses.ArticleRes, error)
	ArticleReadCount(req *requests.ArticleInfo) error
	GetArticleReadCountTopN(req *requests.ArticleTop) ([]responses.ArticleRead, error)
	ArticleEdit(req *requests.ArticleEdit) error
	ArticleCategoryEdit(req *requests.ArticleCategoryEdit) error
	ArticleDel(req *requests.ArticleInfo) error
	GetUserArticleCountTopN(req *requests.ArticleUserTop) ([]responses.ArticleUserCount, error)
}

// DefaultArticleService 默认service所拥有的对象
type DefaultArticleService struct {
	logger     *zap.Logger
	v          *viper.Viper
	pool       *redisPool.Pool
	trace      opentracing.Tracer
	Repository repositories.ArticleRepository
	esConn     *elastic.Client
}

// NewArticleService 初始化ArticleService
func NewArticleService(
	logger *zap.Logger,
	v *viper.Viper,
	pool *redisPool.Pool,
	trace opentracing.Tracer,
	repository repositories.ArticleRepository,
	esConn *elastic.Client,
) ArticleService {
	return &DefaultArticleService{
		logger:     logger.With(zap.String("type", "DefaultArticleService")),
		v:          v,
		pool:       pool,
		trace:      trace,
		Repository: repository,
		esConn:     esConn,
	}
}

// Article 发表文章
func (s *DefaultArticleService) Article(ctx context.Context, req *requests.Article) (responses.Article, error) {
	var result responses.Article
	span, _ := opentracing.StartSpanFromContext(ctx, "article.service")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()
	cateInfo := s.Repository.GetCategory(req.CategoryID)
	if cateInfo.ID == 0 {
		s.logger.Error("categoryId 不存在", zap.Any("category_id", req.CategoryID))
		return result, gorm.ErrRecordNotFound
	}

	res, err := s.Repository.Article(req, ctx)

	if err != nil {
		s.logger.Error("发表文章失败", zap.Error(err))
		return result, errors.New("发表文章失败")
	}
	span.LogFields(
		log.String("event", "insert article service"),
	)

	todayStr := time.Now().Format("20060102")
	keyToday := fmt.Sprintf(redis.KeyUserArticleCount, todayStr)
	keyTotal := fmt.Sprintf(redis.KeyUserArticleCount, constutil.TotalRank)
	keys := make([]string, 0, 2)
	c, err := s.pool.Dial()
	if err != nil || c == nil {
		s.logger.Error("redis连接失败", zap.Error(err))
		return result, errors.New("redis连接失败")
	}

	// 记录用户每天发贴数和总发贴数
	keys = append(keys, keyToday, keyTotal)
	for _, key := range keys {
		_, err = c.Do("ZINCRBY", key, 1, req.UserName)
		if err != nil {
			s.logger.Error("redis操作incr失败", zap.Error(err))
			return result, errors.New("redis操作incr失败")
		}
	}

	result = responses.Article{
		ID:    res.ID,
		Title: res.Title,
	}
	return result, err
}

// GetArticle 获取文章详情
func (s *DefaultArticleService) GetArticle(req *requests.ArticleInfo) (responses.ArticleRes, error) {
	var result responses.ArticleRes
	info := s.Repository.GetArticle(req.ID)
	result = responses.ArticleRes{
		ID:        info.ID,
		Title:     info.Title,
		Summary:   info.Summary,
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
		Category: responses.Category{
			ID:        info.Category.ID,
			Name:      info.Category.Name,
			Num:       info.Category.Num,
			CreatedAt: info.Category.CreatedAt,
			UpdatedAt: info.Category.UpdatedAt,
		},
	}
	return result, nil
}

// ArticleReadCount 文章浏览次数 +1
func (s *DefaultArticleService) ArticleReadCount(req *requests.ArticleInfo) error {
	err := s.Repository.ArticleReadCount(req.ID)
	if err != nil {
		s.logger.Error("文章不存在", zap.Error(err))
		return errors.New("文章不存在")
	}
	todayStr := time.Now().Format("20060102")
	key := fmt.Sprintf(redis.KeyArticleCount, todayStr)
	c, err := s.pool.Dial()
	if err != nil || c == nil {
		s.logger.Error("获取浏览次数失败", zap.Error(err))
		return errors.New("获取浏览次数失败")
	}
	_, err = c.Do("ZINCRBY", key, 1, req.ID)
	if err != nil {
		s.logger.Error("redis操作incr失败", zap.Error(err))
		return errors.New("redis操作incr失败")
	}

	return nil
}

// GetArticleReadCountTopN 通过redis获取TOPN浏览次数的文章
func (s *DefaultArticleService) GetArticleReadCountTopN(req *requests.ArticleTop) ([]responses.ArticleRead, error) {
	var result []responses.ArticleRead
	todayStr := time.Now().Format("20060102")
	key := fmt.Sprintf(redis.KeyArticleCount, todayStr)
	c, err := s.pool.Dial()
	if err != nil || c == nil {
		s.logger.Error("获取浏览次数失败", zap.Error(err))
		return result, errors.New("获取浏览次数失败")
	}
	// zrevrange Key 0 n-1 从redis取出前n位的文章id
	hotArticles, err := redisPool.StringMap(c.Do("ZREVRANGEBYSCORE", key, req.N-1, 0, "withscores"))
	if err != nil {
		s.logger.Error("ZREVRANGEBYSCORE", zap.Any("error", err))
	}
	hotSize := len(hotArticles)
	ids := make([]int, 0, hotSize)
	artMap := make(map[int]int, hotSize) // 保存articleId:count的map
	for idStr, countStr := range hotArticles {
		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			s.logger.Warn("ArticleTopN:strconv.Atoi ID failed", zap.Any("error", err))
			continue
		}
		count, err := strconv.Atoi(countStr)
		if err != nil {
			s.logger.Warn("ArticleTopN:strconv.Atoi count failed", zap.Any("error", err))
			continue
		}
		ids = append(ids, idInt)
		artMap[idInt] = count
	}

	res, err := s.Repository.GetTOPNArticles(ids)
	// 排行榜数据 更新viewCount size如果不设置为0就是默认2个nil
	for _, art := range res {
		count := artMap[art.ID]
		r := responses.ArticleRead{
			ID:    art.ID,
			Title: art.Title,
			Count: count,
		}
		result = append(result, r)
	}

	// 自定义排序 sort.Slice排序 按count排序
	sort.Slice(result, func(i int, j int) bool {
		return result[i].Count > result[j].Count
	})

	return result, nil
}

// ArticleEdit 更新文章  全量更新
func (s *DefaultArticleService) ArticleEdit(req *requests.ArticleEdit) error {
	articleInfo := s.Repository.GetArticle(req.ID)
	cateInfo := s.Repository.GetCategory(req.CategoryID)
	if articleInfo.ID == 0 || cateInfo.ID == 0 {
		s.logger.Error("articleId或categoryId不存在", zap.Any("article_id", req.ID), zap.Any("category_id", req.CategoryID))
		return gorm.ErrRecordNotFound
	}

	article := &models.Article{
		ID:         req.ID,
		CategoryID: req.CategoryID,
		Summary:    req.Summary,
		Title:      req.Title,
		UpdatedAt:  time.Now(),
	}
	return s.Repository.ArticleEdit(article)
}

// ArticleCategoryEdit 更新文章类别  部分更新
func (s *DefaultArticleService) ArticleCategoryEdit(req *requests.ArticleCategoryEdit) error {
	articleInfo := s.Repository.GetArticle(req.ID)
	cateInfo := s.Repository.GetCategory(req.CategoryID)
	if articleInfo.ID == 0 || cateInfo.ID == 0 {
		s.logger.Error("articleId或categoryId不存在", zap.Any("article_id", req.ID), zap.Any("category_id", req.CategoryID))
		return gorm.ErrRecordNotFound
	}

	article := &models.Article{
		ID:         req.ID,
		CategoryID: req.CategoryID,
		UpdatedAt:  time.Now(),
	}
	return s.Repository.ArticleCategoryEdit(article)
}

// ArticleDel 删除文章
func (s *DefaultArticleService) ArticleDel(req *requests.ArticleInfo) error {
	info := s.Repository.GetArticle(req.ID)
	if info.ID == 0 {
		s.logger.Error("articleId不存在", zap.Any("article_id", req.ID))
		return gorm.ErrRecordNotFound
	}
	return s.Repository.ArticleDel(info.ID)
}

// GetUserArticleCountTopN 通过redis获取TOPN用户文章数排行榜
func (s *DefaultArticleService) GetUserArticleCountTopN(req *requests.ArticleUserTop) ([]responses.ArticleUserCount, error) {
	var (
		result  []responses.ArticleUserCount
		keyRank string
	)
	c, err := s.pool.Dial()
	if err != nil || c == nil {
		s.logger.Error("获取用户文章排行榜失败", zap.Error(err))
		return result, errors.New("获取用户文章排行榜失败")
	}

	// 按周排行
	if req.RankType == constutil.WeekRank {
		keyRank = fmt.Sprintf(redis.KeyUserArticleCount, constutil.WeekRank)
		err = redis.UnionStore(constutil.WeekDays, keyRank, c)
	} else if req.RankType == constutil.MonthRank {
		// 当月到今天的天数
		keyRank = fmt.Sprintf(redis.KeyUserArticleCount, constutil.MonthRank)
		err = redis.UnionStore(timeutil.GetNowMonthDays(), keyRank, c)
	} else {
		// 总排行
		keyRank = fmt.Sprintf(redis.KeyUserArticleCount, constutil.TotalRank)
	}
	if err != nil {
		s.logger.Error("合并一周/当月的用户发表文章数", zap.Error(err))
		return result, errors.New("合并一周/当月的用户发表文章数")
	}

	// ZREVRANGEBYSCORE Key 0 n-1 从redis取出前n位的用户username及发贴数
	hotUserArticles, err := redisPool.StringMap(c.Do("ZREVRANGEBYSCORE", keyRank, req.N-1, 0, "withscores"))
	if err != nil {
		s.logger.Error("ZREVRANGEBYSCORE", zap.Any("error", err))
	}
	for username, countStr := range hotUserArticles {
		if err != nil {
			s.logger.Warn("ArticleTopN:strconv.Atoi ID failed", zap.Any("error", err))
			continue
		}
		count, err := strconv.Atoi(countStr)
		if err != nil {
			s.logger.Warn("ArticleTopN:strconv.Atoi count failed", zap.Any("error", err))
			continue
		}
		result = append(result, responses.ArticleUserCount{
			UserName: username,
			Count:    count,
		})
	}

	// 自定义排序 sort.Slice排序 按count排序
	sort.Slice(result, func(i int, j int) bool {
		return result[i].Count > result[j].Count
	})

	return result, nil
}
