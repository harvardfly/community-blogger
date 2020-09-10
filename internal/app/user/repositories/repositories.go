package repositories

import (
	"github.com/google/wire"
)

// ProviderSet user Repository wire
var ProviderSet = wire.NewSet(NewMysqlUserRepository)
