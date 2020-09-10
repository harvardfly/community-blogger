// +build wireinject

package rpc

import (
	"github.com/google/wire"
	"community-blogger/internal/app/userrpc/grpcserver"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/log"
	"community-blogger/internal/app/userrpc/repositories"
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
