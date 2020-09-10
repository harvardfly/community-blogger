package services

import (
	"github.com/google/wire"
)

// ProviderSet home Service wire 注入
var ProviderSet = wire.NewSet(NewHomeService)
