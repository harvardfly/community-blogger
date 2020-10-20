// +build wireinject

package main

import (
	"community-blogger/internal/app/cronjob"
	"community-blogger/internal/app/cronjob/services"
	"community-blogger/internal/pkg/app"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/es"
	"community-blogger/internal/pkg/log"
	"community-blogger/internal/pkg/redis"
	"community-blogger/internal/pkg/transports/cron"
	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	redis.ProviderSet,
	es.ProviderSet,
	cron.ProviderSet,
	cronjob.ProviderSet,
	services.ProviderSet)

func CreateApp(cf string) (*app.Application, error) {
	panic(wire.Build(providerSet))
}
