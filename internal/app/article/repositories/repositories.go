package repositories

import (
	"github.com/google/wire"
)

// ProviderSet article Repository wire
var ProviderSet = wire.NewSet(NewMysqlArticleRepository)
