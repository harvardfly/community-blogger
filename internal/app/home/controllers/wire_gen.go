// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package controllers

import (
	"community-blogger/internal/app/home/repositories"
	"community-blogger/internal/app/home/services"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/log"
	"community-blogger/internal/pkg/storages/minio"
	"github.com/google/wire"
)

// Injectors from wire.go:

func CreateHomeController(cf string, rpo repositories.HomeRepository) (*HomeController, error) {
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
	minioOptions, err := minio.NewOptions(viper, logger)
	if err != nil {
		return nil, err
	}
	client, err := minio.New(minioOptions)
	if err != nil {
		return nil, err
	}
	homeService := services.NewHomeService(logger, viper, rpo, client)
	homeController := NewHomeController(logger, homeService)
	return homeController, nil
}

// wire.go:

var testProviderSet = wire.NewSet(log.ProviderSet, config.ProviderSet, minio.ProviderSet, services.ProviderSet, ProviderSet)
