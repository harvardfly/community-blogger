// +build wireinject

package main

import (
	"community-blogger/internal/pkg/es"
	"community-blogger/internal/pkg/kafka"
	"github.com/google/wire"
	"community-blogger/internal/app/article"
	"community-blogger/internal/app/article/controllers"
	"community-blogger/internal/app/article/repositories"
	"community-blogger/internal/app/article/services"
	"community-blogger/internal/pkg/app"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/database"
	"community-blogger/internal/pkg/jaeger"
	"community-blogger/internal/pkg/log"
	"community-blogger/internal/pkg/redis"
	"community-blogger/internal/pkg/transports/http"
)

var providerSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	database.ProviderSet,
	redis.ProviderSet,
	jaeger.ProviderSet,
	es.ProviderSet,
	kafka.ProviderSet,
	repositories.ProviderSet,
	services.ProviderSet,
	http.ProviderSet,
	article.ProviderSet,
	controllers.ProviderSet)

func CreateApp(cf string) (*app.Application, error) {
	panic(wire.Build(providerSet))
}
