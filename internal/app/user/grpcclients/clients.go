package grpcclients

import (
	"github.com/google/wire"
)

// clientTarget user rpc client settings
type clientTarget struct {
	Caller   string
	Callee   string
	Schema   string
	EtcdAddr string
}

// ProviderSet user rpc client wire
var ProviderSet = wire.NewSet(NewUserClient)
