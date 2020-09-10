// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package grpcserver

import (
	"github.com/google/wire"
	"community-blogger/internal/app/userrpc/repositories"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/log"
)

// Injectors from wire.go:

func CreateUserRpcServer(f string, rpo repositories.UserRepository) (UserRPCServer, error) {
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
	userRPCServer := NewUserRPCServer(logger, viper, rpo)
	return userRPCServer, nil
}

// wire.go:

var testProviderSet = wire.NewSet(log.ProviderSet, config.ProviderSet, ProviderSet)