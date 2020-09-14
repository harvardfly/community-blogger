package middleware

import (
	"community-blogger/internal/pkg/redis"
	"fmt"
	"github.com/gin-gonic/gin"
	redisPool "github.com/gomodule/redigo/redis"
	"net/http"
	"time"
)

// Limiter 限流中间件 用于同一用户发表文章频率限制 避免恶意刷文章
func Limiter(ctx *gin.Context) {
	now := time.Now().UnixNano()
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "username获取失败"})
	}
	key := fmt.Sprintf(redis.KeyLimitArticleUser, username)
	c, err := redis.Client.RedisCon.Dial()
	if err != nil || c == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "redis连接失败"})
		return
	}

	//限制五秒一次请求
	var limit int64 = 1
	dura := time.Second * 5
	//删除有序集合中的五秒之前的数据
	_, err = c.Do("ZREMRANGEBYSCORE", key, "0", fmt.Sprint(now-(dura.Nanoseconds())))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "redis操作ZREMRANGEBYSCORE失败"})
	}
	reqs, _ := redisPool.Int64(c.Do("ZCARD", key))
	if reqs >= limit {
		ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"status":  http.StatusTooManyRequests,
			"message": "too many request",
		})
		return
	}

	ctx.Next()
	_, err = c.Do("ZADD", key, float64(now), float64(now))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "redis操作ZADD失败"})
	}
	_, err = c.Do("EXPIRE", key, dura)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "redis操作EXPIRE失败"})
	}
}

// BucketLimit redis实现令牌桶限流
func BucketLimit(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "username获取失败"})
	}
	key := fmt.Sprintf(redis.KeyBucketLimitArticleUser, username)
	c, err := redis.Client.RedisCon.Dial()
	if err != nil || c == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "redis连接失败"})
		return
	}

	rate := 1                                                       // 令牌生成速度 每秒1个token
	capacity := 1                                                   // 桶容量
	tokens, err := redisPool.Int(c.Do("hget", key, "tokens"))       // 桶中的令牌数
	lastTime, err := redisPool.Int64(c.Do("hget", key, "lastTime")) // 上次令牌生成时间
	now := time.Now().Unix()

	// 初始状态下 令牌数量为桶的容量
	existKey, err := redisPool.Int(c.Do("exists", key))
	if existKey != 1 {
		tokens = capacity
		c.Do("hset", key, "lastTime", now)
	}
	deltaTokens := int(now-lastTime) * rate // 经过一段时间后生成的令牌
	if deltaTokens > 1 {
		tokens = tokens + deltaTokens // 增加令牌
	}
	if tokens > capacity {
		tokens = capacity
		lastTime = time.Now().Unix() // 记录令牌生成时间 秒为单位
		c.Do("hset", key, "lastTime", lastTime)
	}
	if tokens >= 1 {
		tokens -= 1 // 请求进来了，令牌就减少1
		c.Do("hset", key, "tokens", tokens)
		return
	}

	// 无空闲token可用时 429状态码限流提示
	ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
		"status":  http.StatusTooManyRequests,
		"message": "too many request",
	})
}
