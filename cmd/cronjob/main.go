package main

import (
	"flag"
	"fmt"
)

var configFile = flag.String("f", "cronjob.yaml", "set config file which viper will loading.")

func main() {
	flag.Parse()
	fmt.Println(*configFile)
	app, err := CreateApp(*configFile)
	if err != nil {
		panic(err)
	}
	if err := app.Start(); err != nil {
		panic(err)
	}
	app.AwaitSignal()
}
