// +build wireinject

package controllers

import (
	userproto "community-blogger/api/protos/user"
	"community-blogger/internal/app/user/repositories"
	"community-blogger/internal/app/user/services"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/log"
	"github.com/google/wire"
)

var testProviderSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	services.ProviderSet,
	ProviderSet)

func CreateUserController(cf string,
	rpo repositories.UserRepository,
	userclient userproto.UserClient) (*UserController, error) {
	panic(wire.Build(testProviderSet))
}
