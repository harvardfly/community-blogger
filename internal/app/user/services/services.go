package services

import (
	"github.com/google/wire"
)

// ProviderSet 定义user service wire
var ProviderSet = wire.NewSet(NewUserService)
