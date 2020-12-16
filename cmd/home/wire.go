// +build wireinject

package main

import (
	"community-blogger/internal/app/home"
	"community-blogger/internal/app/home/controllers"
	"community-blogger/internal/app/home/repositories"
	"community-blogger/internal/app/home/services"
	"community-blogger/internal/pkg/app"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/database"
	"community-blogger/internal/pkg/log"
	"community-blogger/internal/pkg/storages/minio"
	"community-blogger/internal/pkg/storages/oss"
	"community-blogger/internal/pkg/storages/qiniu"
	"community-blogger/internal/pkg/transports/http"

	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	database.ProviderSet,
	repositories.ProviderSet,
	minio.ProviderSet,
	qiniu.ProviderSet,
	oss.ProviderSet,
	services.ProviderSet,
	http.ProviderSet,
	home.ProviderSet,
	controllers.ProviderSet)

func CreateApp(cf string) (*app.Application, error) {
	panic(wire.Build(providerSet))
}
