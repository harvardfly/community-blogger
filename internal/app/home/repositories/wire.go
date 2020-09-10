// +build wireinject

package repositories

import (
	"github.com/google/wire"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/database"
	"community-blogger/internal/pkg/log"
)

var testProviderSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	database.ProviderSet,
	ProviderSet)

func CreateHomeRepository(f string) (HomeRepository, error) {
	panic(wire.Build(testProviderSet))
}
