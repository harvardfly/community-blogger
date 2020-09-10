package repositories

import (
	"github.com/google/wire"
)

// ProviderSet home Repository wire 注入
var ProviderSet = wire.NewSet(NewMysqlHomeRepository)
