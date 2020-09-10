// +build wireinject

package controllers

import (
	"github.com/google/wire"
	"community-blogger/internal/app/home/repositories"
	"community-blogger/internal/app/home/services"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/log"
)

var testProviderSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	services.ProviderSet,
	ProviderSet)

func CreateHomeController(cf string,
	rpo repositories.HomeRepository) (*HomeController, error) {
	panic(wire.Build(testProviderSet))
}
