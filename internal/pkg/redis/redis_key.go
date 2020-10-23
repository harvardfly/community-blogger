package redis

const (
	// KeyArticleCount 24小时文章阅读数key blog:article:count:20201020
	KeyArticleCount = "blog:article:read:count:%s"
	// KeyLimitArticleUser 限制用户发布文章频率  防止恶意刷文章
	KeyLimitArticleUser = "blog:article:limit:username:%s"
	// KeyBucketLimitArticleUser 令牌桶算法限流
	KeyBucketLimitArticleUser = "blog:article:bucket:limit:username:%s"
	// KeyUserArticleCount 用户发表文章数key blog:user:article:count:20201020
	KeyUserArticleCount = "blog:user:article:count:%s"
)
