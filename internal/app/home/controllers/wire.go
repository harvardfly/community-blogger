// +build wireinject

package controllers

import (
	"community-blogger/internal/app/home/repositories"
	"community-blogger/internal/app/home/services"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/log"
	"community-blogger/internal/pkg/storages/minio"

	"github.com/google/wire"
)

var testProviderSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	minio.ProviderSet,
	services.ProviderSet,
	ProviderSet)

func CreateHomeController(cf string,
	rpo repositories.HomeRepository) (*HomeController, error) {
	panic(wire.Build(testProviderSet))
}
