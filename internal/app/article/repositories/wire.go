// +build wireinject

package repositories

import (
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/database"
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
	database.ProviderSet,
	redis.ProviderSet,
	jaeger.ProviderSet,
	es.ProviderSet,
	kafka.ProviderSet,
	ProviderSet)

func CreateArticleRepository(f string) (ArticleRepository, error) {
	panic(wire.Build(testProviderSet))
}
