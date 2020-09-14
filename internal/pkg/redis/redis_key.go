package redis

const (
	// KeyArticleCount 24小时文章阅读数key eq:blog:article:count:20200506
	KeyArticleCount = "blog:article:read:count:%s"
	// KeyLimitArticleUser 限制用户发布文章频率  防止恶意刷文章
	KeyLimitArticleUser = "blog:article:limit:username:%s"
	// KeyBucketLimitArticleUser 令牌桶算法限流
	KeyBucketLimitArticleUser = "blog:article:bucket:limit:username:%s"
)
