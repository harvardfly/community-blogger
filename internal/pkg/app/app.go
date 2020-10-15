package app

import (
	"community-blogger/internal/pkg/transports/grpc"
	"os"
	"os/signal"
	"syscall"

	"community-blogger/internal/pkg/transports/http"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Application app server
type Application struct {
	name       string
	logger     *zap.Logger
	httpServer *http.Server
	grpcServer *grpc.Server
}

// Option app option
type Option func(*Application) error

// HTTPServerOption app http server option
func HTTPServerOption(svr *http.Server) Option {
	return func(app *Application) error {
		svr.Application(app.name)
		app.httpServer = svr
		return nil
	}
}

// GrpcServerOptions app grpc server option
func GrpcServerOptions(svr *grpc.Server) Option {
	return func(app *Application) error {
		svr.Application(app.name)
		app.grpcServer = svr
		return nil
	}
}

// New new app
func New(name string, logger *zap.Logger, options ...Option) (*Application, error) {
	app := &Application{
		name:   name,
		logger: logger.With(zap.String("type", "Application")),
	}

	for _, option := range options {
		if err := option(app); err != nil {
			return nil, err
		}
	}

	return app, nil
}

// Start start app server
func (a *Application) Start() error {
	if a.httpServer != nil {
		if err := a.httpServer.Start(); err != nil {
			return errors.Wrap(err, "http server start error")
		}
	}

	if a.grpcServer != nil {
		if err := a.grpcServer.Start(); err != nil {
			return errors.Wrap(err, "grpc server start error")
		}
	}

	return nil
}

// AwaitSignal await signal for exit app server
func (a *Application) AwaitSignal() {
	c := make(chan os.Signal, 1)
	signal.Reset(syscall.SIGTERM, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	s := <-c
	a.logger.Info("receive a signal", zap.String("signal", s.String()))
	if a.httpServer != nil {
		if err := a.httpServer.Stop(); err != nil {
			a.logger.Warn("stop http server error", zap.Error(err))
		}
	}
	if a.grpcServer != nil {
		if err := a.grpcServer.Stop(); err != nil {
			a.logger.Warn("stop grpc server error", zap.Error(err))
		}
	}

	os.Exit(0)
}

// ProviderSet wire 注入
var ProviderSet = wire.NewSet(New)
