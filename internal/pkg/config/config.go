package config

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

// New new config
func New(path string) (*viper.Viper, error) {
	var (
		err error
		v   = viper.New()
	)
	v.AddConfigPath(".")
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err == nil {
		fmt.Printf("use config file -> %s\n", v.ConfigFileUsed())
	} else {
		return nil, errors.Wrap(err, "read config file error")
	}

	return v, err
}

// ProviderSet dependency injection
var ProviderSet = wire.NewSet(New)
