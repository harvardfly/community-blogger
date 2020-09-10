// +build wireinject

package main

import (
	"github.com/google/wire"
	"community-blogger/internal/app/home"
	"community-blogger/internal/app/home/controllers"
	"community-blogger/internal/app/home/repositories"
	"community-blogger/internal/app/home/services"
	"community-blogger/internal/pkg/app"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/database"
	"community-blogger/internal/pkg/log"
	"community-blogger/internal/pkg/transports/http"
)

var providerSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	database.ProviderSet,
	repositories.ProviderSet,
	services.ProviderSet,
	http.ProviderSet,
	home.ProviderSet,
	controllers.ProviderSet)

func CreateApp(cf string) (*app.Application, error) {
	panic(wire.Build(providerSet))
}
