// +build wireinject

package services

import (
	"github.com/google/wire"
	"community-blogger/internal/app/article/repositories"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/jaeger"
	"community-blogger/internal/pkg/log"
	"community-blogger/internal/pkg/redis"
)

var testProviderSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	redis.ProviderSet,
	jaeger.ProviderSet,
	ProviderSet)

func CreateArticleService(cf string,
	rpo repositories.ArticleRepository,
) (ArticleService, error) {
	panic(wire.Build(testProviderSet))
}
