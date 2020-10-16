package rpc

import (
	userproto "community-blogger/api/protos/user"
	"community-blogger/internal/pkg/transports/grpc"
	"github.com/google/wire"
	grpc2 "google.golang.org/grpc"
)

// CreateInitServersFn RPC服务入口
func CreateInitServersFn(userSer *UserSer) grpc.InitServers {
	return func(server *grpc2.Server) {
		userproto.RegisterUserServer(server, userSer)
	}
}

// ProviderSet RPC服务依赖注入
var ProviderSet = wire.NewSet(NewUserSer, CreateInitServersFn)
