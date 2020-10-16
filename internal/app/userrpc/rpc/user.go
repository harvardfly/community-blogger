package rpc

import (
	userpb "community-blogger/api/protos/user"
	"community-blogger/internal/app/userrpc/grpcserver"
	"context"
	"go.uber.org/zap"
)

// UserSer 定义userRPC模块UserSer
type UserSer struct {
	logger  *zap.Logger
	service grpcserver.UserRPCServer
}

// NewUserSer 初始化用户RPC模块server
func NewUserSer(logger *zap.Logger, s grpcserver.UserRPCServer) *UserSer {
	return &UserSer{
		logger:  logger.With(zap.String("type", "UserSer")),
		service: s,
	}
}

// FindByToken 实现根据Token获取用户信息RPC方法
func (s *UserSer) FindByToken(ctx context.Context, req *userpb.FindByTokenRequest) (*userpb.UserResponse, error) {
	member, err := s.service.FindByToken(ctx, req)
	return member, err
}

// FindById  实现根据ID获取用户信息RPC方法
func (s *UserSer) FindById(ctx context.Context, req *userpb.FindByIDRequest) (*userpb.UserResponse, error) {
	member, err := s.service.FindByID(ctx, req)
	if err != nil {
		s.logger.Error("user rpc 请求错误", zap.Error(err))
	}
	return member, err
}
