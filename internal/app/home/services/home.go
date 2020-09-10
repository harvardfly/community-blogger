package services

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"community-blogger/internal/app/home/repositories"
	"community-blogger/internal/pkg/requests"
	"community-blogger/internal/pkg/responses"
	"community-blogger/internal/pkg/utils/constutil"
)

// HomeService home模块service定义
type HomeService interface {
	HomeList(req *requests.HomeList) (int, []responses.Home, error)
	Home(req *requests.Home) error
}

// DefaultHomeService home模块service默认对象
type DefaultHomeService struct {
	logger     *zap.Logger
	v          *viper.Viper
	Repository repositories.HomeRepository
}

// NewHomeService 初始化
func NewHomeService(
	logger *zap.Logger,
	v *viper.Viper,
	repository repositories.HomeRepository,
) HomeService {
	return &DefaultHomeService{
		logger:     logger.With(zap.String("type", "DefaultHomeService")),
		v:          v,
		Repository: repository,
	}
}

// HomeList home信息列表
func (s *DefaultHomeService) HomeList(req *requests.HomeList) (int, []responses.Home, error) {
	if req.Page < 0 {
		req.Page = 0
	}
	if req.Limit <= 0 {
		req.Limit = constutil.DefaultPageSize
	} else if req.Limit > constutil.MaxPageSize {
		req.Limit = constutil.MaxPageSize
	}
	req.Paginator = map[string]int{"limit": req.Limit, "offset": req.Page}
	req.Orders = map[string]string{"order": req.Order, "desc": req.Description}
	return s.Repository.HomeList(req)
}

// Home 新建home信息
func (s *DefaultHomeService) Home(req *requests.Home) error {
	return s.Repository.Home(req)
}
