package grpc

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"community-blogger/internal/pkg/etcdservice"
	"community-blogger/internal/pkg/utils/netutil"
)

// ServerOptions grpc server option
type ServerOptions struct {
	Port        int
	Schema      string
	EtcdAddr    string
	ServiceName string
	TTL         int64
}

// NewServerOptions grpc new option
func NewServerOptions(v *viper.Viper, logger *zap.Logger) (*ServerOptions, error) {
	var (
		err error
		o   = new(ServerOptions)
	)

	if err = v.UnmarshalKey("grpc", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal grpc server option error")
	}

	logger.Info("load grpc options success", zap.Any("grpc options", o))

	return o, nil
}

// Server grpc Server
type Server struct {
	o      *ServerOptions
	app    string
	host   string
	port   int
	logger *zap.Logger
	server *grpc.Server
}

// InitServers 初始化Servers
type InitServers func(*grpc.Server)

// NewServer initialize grpc server
func NewServer(o *ServerOptions, logger *zap.Logger, init InitServers) (*Server, error) {
	var gs *grpc.Server
	logger = logger.With(zap.String("type", "grpc"))
	gs = grpc.NewServer()
	init(gs)

	return &Server{
		o:      o,
		logger: logger.With(zap.String("type", "grpc.Server")),
		server: gs,
	}, nil
}

// Application 服务应用
func (s *Server) Application(name string) {
	s.app = name
}

// Start 启动RPC服务
func (s *Server) Start() error {
	s.port = s.o.Port
	if s.port == 0 {
		s.port = netutil.GetAvailablePort()
	}

	addr := fmt.Sprintf(":%d", s.port)

	s.logger.Info("grpc server starting ...", zap.String("addr", addr))

	//将服务地址注册到etcd中
	go etcdservice.Register(s.o.EtcdAddr, s.o.ServiceName, addr, s.o.TTL)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		sig := <-ch
		etcdservice.UnRegister(s.o.ServiceName, addr)

		if i, ok := sig.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}

	}()

	go func() {
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			s.logger.Fatal("failed to listen: %v", zap.Error(err))
		}

		if err := s.server.Serve(lis); err != nil {
			s.logger.Fatal("failed to serve: %v", zap.Error(err))
		}
	}()

	return nil
}

// Stop  停止RPC服务
func (s *Server) Stop() error {
	s.logger.Info("grpc server stopping ...")
	s.server.GracefulStop()
	return nil
}
