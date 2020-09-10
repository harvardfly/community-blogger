// +build wireinject

package services

import (
	"github.com/google/wire"
	"community-blogger/internal/app/home/repositories"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/log"
)

var testProviderSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	ProviderSet)

func CreateHomeService(cf string,
	rpo repositories.HomeRepository,
) (HomeService, error) {
	panic(wire.Build(testProviderSet))
}
