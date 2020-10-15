package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"community-blogger/internal/pkg/utils/netutil"
	"github.com/google/wire"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	ginZap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Options http option
type Options struct {
	Host string
	Port int
	Mode string
}

// Server http server
type Server struct {
	o          *Options
	app        string
	host       string
	port       int
	logger     *zap.Logger
	router     *gin.Engine
	httpServer http.Server
}

// NewOptions new http option
func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("http", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal http option error")
	}

	logger.Info("load http options success", zap.Any("http options", o))
	return o, err
}

// InitControllers init controllers
type InitControllers func(r *gin.Engine)

// NewRouter new gin router
// init gin configuration
func NewRouter(o *Options, logger *zap.Logger, init InitControllers) *gin.Engine {
	// 配置gin
	gin.SetMode(o.Mode)
	r := gin.New()
	// panic之后自动恢复
	r.Use(gin.Recovery())
	// 日志格式化
	r.Use(ginZap.Ginzap(logger, time.RFC3339, true))
	// panic日志格式化
	r.Use(ginZap.RecoveryWithZap(logger, true))
	init(r)

	return r
}

// New new http server
func New(o *Options, logger *zap.Logger, router *gin.Engine) (*Server, error) {
	return &Server{
		o:      o,
		logger: logger,
		router: router,
	}, nil
}

// Application set app name
func (s *Server) Application(name string) {
	s.app = name
}

// Start start app and register to consul
func (s *Server) Start() error {
	s.port = s.o.Port
	if s.port == 0 {
		s.port = netutil.GetAvailablePort()
	}

	s.host = s.o.Host
	if s.host == "" {
		s.host = "127.0.0.1"
	}

	s.httpServer = http.Server{Addr: fmt.Sprintf("%s:%d", s.host, s.port), Handler: s.router}
	s.logger.Info("http server starting...", zap.String("addr", s.httpServer.Addr))

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("start http server error", zap.Error(err))
			return
		}
	}()
	return nil
}

// Stop stop app
func (s *Server) Stop() error {
	s.logger.Info("stopping http server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) // 平滑关闭,等待5秒钟处理
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutdown http server error")
	}

	return nil
}

// ProviderSet dependency injection
var ProviderSet = wire.NewSet(New, NewRouter, NewOptions)
