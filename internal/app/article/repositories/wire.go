// +build wireinject

package repositories

import (
	"github.com/google/wire"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/database"
	"community-blogger/internal/pkg/jaeger"
	"community-blogger/internal/pkg/log"
	"community-blogger/internal/pkg/redis"
)

var testProviderSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	database.ProviderSet,
	redis.ProviderSet,
	jaeger.ProviderSet,
	ProviderSet)

func CreateArticleRepository(f string) (ArticleRepository, error) {
	panic(wire.Build(testProviderSet))
}
