package cron

import (
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/rfyiamcool/cronlib"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Options struct {
	Projects       map[string]string
	EnableAsync    bool
	EnableTryCatch bool
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("cron", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal cron option error")
	}
	logger.Info("load cron options success", zap.Any("cron options", o))
	return o, err
}

type ServerOptional struct {
	spec string
	f    func()
}

type Server struct {
	app    string
	o      *Options
	logger *zap.Logger
	cron   *cronlib.CronSchduler
	jobs   map[string]ServerOptional
}

type InitServers map[string]func()

func New(o *Options, logger *zap.Logger, init InitServers) (*Server, error) {
	optionals := make(map[string]ServerOptional)
	for name, spec := range o.Projects {
		if jobFunc, ok := init[name]; ok {
			optionals[name] = ServerOptional{
				spec: spec,
				f:    jobFunc,
			}
		} else {
			logger.Error("定时任务不存在", zap.String("name", name))
			return nil, errors.New("定时任务不存在")
		}
	}
	return &Server{
		o:      o,
		logger: logger.With(zap.String("type", "cronServer")),
		cron:   cronlib.New(),
		jobs:   optionals,
	}, nil
}

func (s *Server) Application(name string) {
	s.app = name
}

func (s *Server) Start() error {
	go func() {
		if err := s.register(); err != nil {
			s.logger.Fatal("failed to register cron: %v", zap.Error(err))
		}
		s.cron.Start()
		s.cron.Wait()
	}()
	return nil
}

func (s *Server) Register(name string) error {
	job, ok := s.jobs[name]
	if !ok {
		s.logger.Error("定时任务不存在", zap.String("name", name))
		return errors.New("定时任务不存在")
	}
	m, err := cronlib.NewJobModel(job.spec, job.f)
	if err != nil {
		s.logger.Error("创建job失败", zap.Error(err))
		return err
	}
	m.SetAsyncMode(s.o.EnableAsync)
	m.SetTryCatch(s.o.EnableTryCatch)
	return s.cron.UpdateJobModel(name, m)
}

func (s *Server) register() error {
	for name, job := range s.jobs {
		m, err := cronlib.NewJobModel(job.spec, job.f)
		if err != nil {
			s.logger.Error("创建job失败", zap.Error(err))
			return err
		}
		m.SetAsyncMode(s.o.EnableAsync)
		m.SetTryCatch(s.o.EnableTryCatch)
		err = s.cron.Register(name, m)
		if err != nil {
			s.logger.Error("注册job失败", zap.Error(err))
			return err
		}
		s.logger.Info("注册cron任务成功", zap.String("name", name))
	}
	return nil
}

func (s *Server) DeRegister(name string) error {
	if _, ok := s.jobs[name]; !ok {
		s.logger.Error("定时任务不存在", zap.String("name", name))
		return errors.New("定时任务不存在")
	}
	s.cron.StopService(name)
	return nil
}

func (s *Server) deRegister() error {
	for name := range s.jobs {
		s.cron.StopService(name)
		s.logger.Info("deregister cron service success", zap.String("name", name))
	}
	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("cron server stopping ...")
	if err := s.deRegister(); err != nil {
		return errors.Wrap(err, "deregister cron server error")
	}
	return nil
}

var ProviderSet = wire.NewSet(New, NewOptions)
