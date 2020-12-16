// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"community-blogger/internal/app/home"
	"community-blogger/internal/app/home/controllers"
	"community-blogger/internal/app/home/repositories"
	"community-blogger/internal/app/home/services"
	"community-blogger/internal/pkg/app"
	"community-blogger/internal/pkg/config"
	"community-blogger/internal/pkg/database"
	"community-blogger/internal/pkg/log"
	"community-blogger/internal/pkg/storages/minio"
	"community-blogger/internal/pkg/storages/oss"
	"community-blogger/internal/pkg/storages/qiniu"
	"community-blogger/internal/pkg/transports/http"
	"github.com/google/wire"
)

// Injectors from wire.go:

func CreateApp(cf string) (*app.Application, error) {
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
	homeOptions, err := home.NewOptions(viper, logger)
	if err != nil {
		return nil, err
	}
	httpOptions, err := http.NewOptions(viper, logger)
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
	homeRepository := repositories.NewMysqlHomeRepository(logger, databaseDatabase)
	minioOptions, err := minio.NewOptions(viper, logger)
	if err != nil {
		return nil, err
	}
	client, err := minio.New(minioOptions)
	if err != nil {
		return nil, err
	}
	qiniuOptions, err := qiniu.NewOptions(viper, logger)
	if err != nil {
		return nil, err
	}
	formUploader, err := qiniu.New(qiniuOptions)
	if err != nil {
		return nil, err
	}
	ossOptions, err := oss.NewOptions(viper, logger)
	if err != nil {
		return nil, err
	}
	ossClient, err := oss.New(ossOptions)
	if err != nil {
		return nil, err
	}
	homeService := services.NewHomeService(logger, viper, homeRepository, client, formUploader, ossClient)
	homeController := controllers.NewHomeController(logger, homeService)
	initControllers := controllers.CreateInitControllersFn(homeController)
	engine := http.NewRouter(httpOptions, logger, initControllers)
	server, err := http.New(httpOptions, logger, engine)
	if err != nil {
		return nil, err
	}
	application, err := home.NewApp(homeOptions, logger, server)
	if err != nil {
		return nil, err
	}
	return application, nil
}

// wire.go:

var providerSet = wire.NewSet(log.ProviderSet, config.ProviderSet, database.ProviderSet, repositories.ProviderSet, minio.ProviderSet, qiniu.ProviderSet, oss.ProviderSet, services.ProviderSet, http.ProviderSet, home.ProviderSet, controllers.ProviderSet)
