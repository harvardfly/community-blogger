package services

import (
	"github.com/google/wire"
)

// ProviderSet 定义article service wire
var ProviderSet = wire.NewSet(NewArticleService)
