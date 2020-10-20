package services

import (
	"community-blogger/internal/pkg/es"
	"community-blogger/internal/pkg/redis"
	"errors"
	"fmt"
	redisPool "github.com/gomodule/redigo/redis"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// DefaultCronJobService CronJob模块service默认对象
type DefaultCronJobService struct {
	logger  *zap.Logger
	v       *viper.Viper
	pool    *redisPool.Pool
	esConn  *elastic.Client
}

// NewDefaultCronJobService 初始化
func NewDefaultCronJobService(
	logger *zap.Logger,
	v *viper.Viper,
	pool *redisPool.Pool,
	esConn *elastic.Client,
) *DefaultCronJobService {
	return &DefaultCronJobService{
		logger:  logger.With(zap.String("type", "DefaultCronJobService")),
		v:       v,
		pool:    pool,
		esConn:  esConn,
	}
}

// RedisToES 同步redis数据到ES中 数据持久化
func (s *DefaultCronJobService) RedisToES() error {
	todayStr := time.Now().Format("20060102")
	key := fmt.Sprintf(redis.KeyArticleCount, todayStr)
	c, err := s.pool.Dial()
	if err != nil || c == nil {
		s.logger.Error("获取浏览次数失败", zap.Error(err))
		return errors.New("获取浏览次数失败")
	}

	// zrange Key 0 -1 从redis取出zset每天所有的文章id及浏览次数
	articles, err := redisPool.StringMap(c.Do("ZRANGE", key, 0, -1, "withscores"))
	if err != nil {
		s.logger.Error("ZRANGE", zap.Any("error", err))
		return err
	}
	for idStr, countStr := range articles {
		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			s.logger.Warn("Article:strconv.Atoi ID failed", zap.Any("error", err))
			continue
		}
		count, err := strconv.Atoi(countStr)
		if err != nil {
			s.logger.Warn("Article:strconv.Atoi count failed", zap.Any("error", err))
			continue
		}

		// 更新es的文章浏览次数
		Params := make(map[string]string)
		Doc := make(map[string]string)
		Params["index"] = "article"
		Params["id"] = strconv.Itoa(idInt)
		Doc["count"] = strconv.Itoa(count)

		_ = es.Client.Update(Params, Doc)
	}
	return nil
}
