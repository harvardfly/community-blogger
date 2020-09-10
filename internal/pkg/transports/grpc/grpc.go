package grpc

import "github.com/google/wire"

// ProviderSet grpc传输方式 依赖注入
var ProviderSet = wire.NewSet(NewServerOptions, NewServer, NewClientOptions, NewClient)
