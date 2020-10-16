// +build wireinject

package grpcserver

import (
	"community-blogger/internal/app/userrpc/repositories"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/log"
	"github.com/google/wire"
)

var testProviderSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	ProviderSet)

func CreateUserRpcServer(f string,
	rpo repositories.UserRepository) (UserRPCServer, error) {
	panic(wire.Build(testProviderSet))
}
