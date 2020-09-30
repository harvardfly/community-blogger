package main

import (
	"community-blogger/internal/pkg/kafka"
	"flag"
)

var configFile = flag.String("f", "article.yaml", "set config file which viper will loading.")

func main() {
	flag.Parse()
	app, err := CreateApp(*configFile)
	if err != nil {
		panic(err)
	}
	if err := app.Start(); err != nil {
		panic(err)
	}
	// kafka consumer
	go kafka.Client.Subscribe()
	app.AwaitSignal()
}
