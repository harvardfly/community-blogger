// +build wireinject

package services

import (
	"github.com/google/wire"
	userproto "community-blogger/api/protos/user"
	"community-blogger/internal/app/user/repositories"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/database"
	"community-blogger/internal/pkg/log"
)

var testProviderSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	database.ProviderSet,
	ProviderSet)

func CreateUserService(cf string,
	rpo repositories.UserRepository,
	userclient userproto.UserClient) (UserService, error) {
	panic(wire.Build(testProviderSet))
}
