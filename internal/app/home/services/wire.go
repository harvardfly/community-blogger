// +build wireinject

package services

import (
	"community-blogger/internal/app/home/repositories"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/log"
	"community-blogger/internal/pkg/storages/minio"

	"github.com/google/wire"
)

var testProviderSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	minio.ProviderSet,
	ProviderSet)

func CreateHomeService(cf string,
	rpo repositories.HomeRepository,
) (HomeService, error) {
	panic(wire.Build(testProviderSet))
}
