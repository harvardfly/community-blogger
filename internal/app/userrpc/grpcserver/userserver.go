package grpcserver

import (
	userpb "community-blogger/api/protos/user"
	"community-blogger/internal/app/userrpc/repositories"
	"context"
	"errors"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

/*
实现rpc的方法：
	FindByToken(context.Context, *FindByTokenRequest) (*UserResponse, error)
	FindById(context.Context, *FindByIdRequest) (*UserResponse, error)
*/

// UserRPCServer 定义userRpc service
type UserRPCServer interface {
	FindByToken(ctx context.Context, req *userpb.FindByTokenRequest) (*userpb.UserResponse, error)
	FindByID(ctx context.Context, req *userpb.FindByIDRequest) (*userpb.UserResponse, error)
}

var (
	// ErrNotFound 定义错误
	ErrNotFound = errors.New("用户不存在")
)

// DefaultUserService 默认service所拥有的对象
type DefaultUserService struct {
	logger     *zap.Logger
	v          *viper.Viper
	Repository repositories.UserRepository
}

// NewUserRPCServer 初始化UserService
func NewUserRPCServer(
	logger *zap.Logger,
	v *viper.Viper,
	repository repositories.UserRepository,
) UserRPCServer {
	return &DefaultUserService{
		logger:     logger.With(zap.String("type", "DefaultUserService")),
		v:          v,
		Repository: repository,
	}
}

// FindByToken 实现rpc服务通过Token获取用户信息
func (s *DefaultUserService) FindByToken(ctx context.Context, req *userpb.FindByTokenRequest) (*userpb.UserResponse, error) {
	member, err := s.Repository.FindByToken(req.Token)
	if err != nil {
		return &userpb.UserResponse{}, ErrNotFound
	}
	return &userpb.UserResponse{
		Id:       int32(member.ID),
		Token:    member.Token,
		Username: member.Username,
		Password: member.Password,
	}, nil
}

// FindByID 调用rpc服务通过ID获取用户信息
func (s *DefaultUserService) FindByID(ctx context.Context, req *userpb.FindByIDRequest) (*userpb.UserResponse, error) {
	member, err := s.Repository.FindByID(req.Id)
	if err != nil {
		return &userpb.UserResponse{}, ErrNotFound
	}
	return &userpb.UserResponse{
		Id:       int32(member.ID),
		Token:    member.Token,
		Username: member.Username,
		Password: member.Password,
	}, nil
}
