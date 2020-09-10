package grpcserver

import (
	"github.com/google/wire"
)

// ProviderSet 定义grpc service wire
var ProviderSet = wire.NewSet(NewUserRPCServer)
