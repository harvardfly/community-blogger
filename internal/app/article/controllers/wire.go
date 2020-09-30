// +build wireinject

package controllers

import (
	"community-blogger/internal/app/article/repositories"
	"community-blogger/internal/app/article/services"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/es"
	"community-blogger/internal/pkg/jaeger"
	"community-blogger/internal/pkg/kafka"
	"community-blogger/internal/pkg/log"
	"community-blogger/internal/pkg/redis"
	"github.com/google/wire"
)

var testProviderSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	redis.ProviderSet,
	jaeger.ProviderSet,
	es.ProviderSet,
	kafka.ProviderSet,
	services.ProviderSet,
	ProviderSet)

func CreateArticleController(cf string,
	rpo repositories.ArticleRepository) (*ArticleController, error) {
	panic(wire.Build(testProviderSet))
}
