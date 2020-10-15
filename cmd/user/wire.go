// +build wireinject

package main

import (
	"community-blogger/internal/app/user"
	"community-blogger/internal/app/user/controllers"
	"community-blogger/internal/app/user/grpcclients"
	"community-blogger/internal/app/user/repositories"
	"community-blogger/internal/app/user/services"
	"community-blogger/internal/pkg/app"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/database"
	"community-blogger/internal/pkg/log"
	"community-blogger/internal/pkg/transports/grpc"
	"community-blogger/internal/pkg/transports/http"
	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	database.ProviderSet,
	repositories.ProviderSet,
	services.ProviderSet,
	grpc.ProviderSet,
	grpcclients.ProviderSet,
	http.ProviderSet,
	user.ProviderSet,
	controllers.ProviderSet)

func CreateApp(cf string) (*app.Application, error) {
	panic(wire.Build(providerSet))
}
