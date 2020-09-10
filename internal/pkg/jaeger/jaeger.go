package jaeger

import (
	"fmt"
	"github.com/google/wire"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"io"
)

// ClientType 定义jaeger client 结构体
type ClientType struct {
	Tracer opentracing.Tracer
	Closer io.Closer
}

// Client  jaeger连接类型
var Client ClientType

// Options jaeger option
type Options struct {
	Type               string  // const
	Param              float64 // 1
	LogSpans           bool    // true
	LocalAgentHostPort string  // host:port
	Service            string  // service name
}

// NewOptions for jaeger
func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("jaeger", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal redis option error")
	}

	logger.Info("load jaeger options success", zap.Any("jaeger options", o))
	return o, err
}

// New returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func New(o *Options) (opentracing.Tracer, error) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  o.Type,
			Param: o.Param,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: o.LogSpans,
			// 注意：填下地址不能加http://
			LocalAgentHostPort: o.LocalAgentHostPort,
		},
	}
	tracer, closer, err := cfg.New(o.Service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	Client.Tracer = tracer
	Client.Closer = closer

	return tracer, err
}

// ProviderSet inject jaeger settings
var ProviderSet = wire.NewSet(New, NewOptions)
