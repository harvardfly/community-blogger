package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// ClientType 定义redis client 结构体
type ClientType struct {
	RedisCon *redis.Pool
}

// Client  redis连接类型
var Client ClientType

// Options redis option
type Options struct {
	URL         string // host:port
	MaxIdle     int    // 最大空闲连接数
	MaxActive   int    // 最大连接数
	IdleTimeout int    // 空闲连接超时时间
	Timeout     int    // 操作超时时间
	Network     string // tcp or udp
	Password    string
}

// NewOptions for redis
func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("redis", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal redis option error")
	}

	logger.Info("load redis options success", zap.Any("redis options", o))
	return o, err
}

// New redis pool conn
func New(o *Options) (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxIdle:     o.MaxIdle,
		MaxActive:   o.MaxActive,
		IdleTimeout: time.Duration(o.IdleTimeout) * time.Second,
		Wait:        true,
		// Other pool configuration not shown in this example.
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", o.URL,
				redis.DialPassword(o.Password),
				redis.DialConnectTimeout(time.Duration(o.Timeout)*time.Second),
				redis.DialReadTimeout(time.Duration(o.Timeout)*time.Second),
				redis.DialWriteTimeout(time.Duration(o.Timeout)*time.Second))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}
	Client.RedisCon = pool
	return pool, nil
}

// UnionStore 合并zset的key
func UnionStore(rankDays int, keyRank string, c redis.Conn) error {
	today := time.Now()
	unionKeys := make([]interface{}, 0, rankDays+3)
	unionKeys = append(unionKeys, keyRank, rankDays)
	for i := 0; i < rankDays; i++ {
		key := fmt.Sprintf(KeyUserArticleCount, today.AddDate(0, 0, -i).Format("20060102"))
		unionKeys = append(unionKeys, key)
	}

	// 合并一周/当月的用户发表文章数
	_, err := c.Do("ZUNIONSTORE", unionKeys...)
	if err != nil {
		return err
	}
	return nil
}

// ProviderSet inject redis settings
var ProviderSet = wire.NewSet(New, NewOptions)
