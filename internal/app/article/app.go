package article

import (
	"community-blogger/internal/pkg/app"
	"community-blogger/internal/pkg/transports/http"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Options 定义文章类配置选项
type Options struct {
	Name string
}

// NewOptions 定义文章类初始化方法
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
func NewApp(o *Options, logger *zap.Logger, hs *http.Server) (*app.Application, error) {
	a, err := app.New(o.Name, logger, app.HTTPServerOption(hs))
	if err != nil {
		return nil, errors.Wrap(err, "new app error")
	}

	return a, nil
}

// ProviderSet article模块wire NewSet
var ProviderSet = wire.NewSet(NewApp, NewOptions)
