// +build wireinject

package main

import (
	"community-blogger/internal/app/userrpc"
	"community-blogger/internal/app/userrpc/grpcserver"
	"community-blogger/internal/app/userrpc/repositories"
	"community-blogger/internal/app/userrpc/rpc"
	"community-blogger/internal/pkg/app"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/database"
	"community-blogger/internal/pkg/log"
	"community-blogger/internal/pkg/transports/grpc"
	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	database.ProviderSet,
	repositories.ProviderSet,
	grpcserver.ProviderSet,
	grpc.ProviderSet,
	userrpc.ProviderSet,
	rpc.ProviderSet,
)

func CreateApp(cf string) (*app.Application, error) {
	panic(wire.Build(providerSet))
}
