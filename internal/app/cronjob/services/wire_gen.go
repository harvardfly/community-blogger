// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package services

import (
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/es"
	"community-blogger/internal/pkg/log"
	"community-blogger/internal/pkg/redis"
	"community-blogger/internal/pkg/transports/cron"
	"github.com/google/wire"
)

// Injectors from wire.go:

func CreateDefaultCronJobService(cf string) (*DefaultCronJobService, error) {
	viper, err := config.New(cf)
	if err != nil {
		return nil, err
	}
	options, err := log.NewOptions(viper)
	if err != nil {
		return nil, err
	}
	logger, err := log.New(options)
	if err != nil {
		return nil, err
	}
	redisOptions, err := redis.NewOptions(viper, logger)
	if err != nil {
		return nil, err
	}
	pool, err := redis.New(redisOptions)
	if err != nil {
		return nil, err
	}
	esOptions := es.NewOptions(viper, logger)
	client := es.New(esOptions, logger)
	defaultCronJobService := NewDefaultCronJobService(logger, viper, pool, client)
	return defaultCronJobService, nil
}

// wire.go:

var testProviderSet = wire.NewSet(log.ProviderSet, config.ProviderSet, redis.ProviderSet, es.ProviderSet, cron.ProviderSet, ProviderSet)