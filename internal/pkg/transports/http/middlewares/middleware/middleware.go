package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	redisPool "github.com/gomodule/redigo/redis"
	"log"
	"net/http"
	"time"
	"community-blogger/internal/pkg/database"
	"community-blogger/internal/pkg/models"
	"community-blogger/internal/pkg/redis"
)

var Routers = map[string]map[string]string{
	"post_article": {
		"method": "POST",
		"path":   "/api/v1/article",
	},
}

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

// AccumulatePoints 积分中间件 用于同一用户发表文章增加积分的规则
func AccumulatePoints(ctx *gin.Context) {
	method := ctx.Request.Method   // POST
	path := ctx.Request.RequestURI // /api/v1/article
	action := fmt.Sprintf("%s-%s", method, path)
	ctx.Next()
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "username获取失败"})
	}
	key := fmt.Sprintf(redis.KeyArticlePostPoints, username, action)
	c, err := redis.Client.RedisCon.Dial()
	if err != nil || c == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "redis连接失败"})
		return
	}
	points, _ := redisPool.Int64(c.Do("GET", key))
	switch action {
	case fmt.Sprintf("%s-%s", Routers["method"], Routers["path"]):
		//限制20次
		var limit int64 = 20
		expire := time.Second * 60 * 60 * 24
		if points >= limit {
			// 今日没有积分奖励了
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  http.StatusTooManyRequests,
				"message": "too many request",
			})
			return
		}
		points, err = redisPool.Int64(c.Do("INCRBY", key, 1))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "redis操作ZADD失败"})
		}
		_, err = c.Do("EXPIRE", key, expire)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "redis操作EXPIRE失败"})
		}
	}
	// 积分对应等级
	level := PointToLevel(points)
	fmt.Println(level)
	// 积分、等级入库更新
	go func() {
		err = database.DBClient.Mysql.Model(&models.Integral{}).Create(&models.Integral{}).Error
		if err != nil {
			log.Fatal()
			return
		}
	}()
}

func PointToLevel(point int64) string {
	return ""
}
