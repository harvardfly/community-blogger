// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package repositories

import (
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/database"
	"community-blogger/internal/pkg/jaeger"
	"community-blogger/internal/pkg/log"
	"community-blogger/internal/pkg/redis"
	"github.com/google/wire"
)

// Injectors from wire.go:

func CreateArticleRepository(f string) (ArticleRepository, error) {
	viper, err := config.New(f)
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
	databaseOptions, err := database.NewOptions(viper, logger)
	if err != nil {
		return nil, err
	}
	databaseDatabase, err := database.New(databaseOptions)
	if err != nil {
		return nil, err
	}
	articleRepository := NewMysqlArticleRepository(logger, databaseDatabase)
	return articleRepository, nil
}

// wire.go:

var testProviderSet = wire.NewSet(log.ProviderSet, config.ProviderSet, database.ProviderSet, redis.ProviderSet, jaeger.ProviderSet, ProviderSet)