// +build wireinject

package services

import (
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/es"
	"community-blogger/internal/pkg/log"
	"community-blogger/internal/pkg/redis"
	"community-blogger/internal/pkg/transports/cron"
	"github.com/google/wire"
)

var testProviderSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	redis.ProviderSet,
	es.ProviderSet,
	cron.ProviderSet,
	ProviderSet)

func CreateDefaultCronJobService(cf string,
) (*DefaultCronJobService, error) {
	panic(wire.Build(testProviderSet))
}
