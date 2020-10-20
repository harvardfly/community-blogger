package cronjob

import (
	"community-blogger/internal/pkg/app"
	"community-blogger/internal/pkg/transports/cron"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Options 定义cron模块 Options
type Options struct {
	Name string
}

// NewOptions 初始化Options
func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("app", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal app option error")
	}

	logger.Info("load application options success")
	return o, err
}

// NewApp 初始化app
func NewApp(o *Options, logger *zap.Logger, con *cron.Server) (*app.Application, error) {
	a, err := app.New(o.Name, logger, app.CronServerOptions(con))
	if err != nil {
		return nil, errors.Wrap(err, "new app error")
	}

	return a, nil
}

// ProviderSet 定义wire provider
var ProviderSet = wire.NewSet(NewApp, NewOptions)
