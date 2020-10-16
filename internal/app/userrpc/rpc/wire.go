// +build wireinject

package rpc

import (
	"community-blogger/internal/app/userrpc/grpcserver"
	"community-blogger/internal/app/userrpc/repositories"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/log"
	"github.com/google/wire"
)

var testProviderSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	grpcserver.ProviderSet,
	ProviderSet)

func CreateUserRpc(cf string,
	rpo repositories.UserRepository) (*UserSer, error) {
	panic(wire.Build(testProviderSet))
}
