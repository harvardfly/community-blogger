package grpcclients

import (
	"github.com/google/wire"
)

// clientTarget user rpc client settings
type clientTarget struct {
	User     string
	Caller   string
	Callee   string
	EtcdAddr string
}

// ProviderSet user rpc client wire
var ProviderSet = wire.NewSet(NewUserClient)
