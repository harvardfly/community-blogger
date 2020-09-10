package repositories

import (
	"github.com/google/wire"
)

// ProviderSet userRpc Repository wire
var ProviderSet = wire.NewSet(NewMysqlUserRepository)
