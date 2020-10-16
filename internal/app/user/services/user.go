package services

import (
	userproto "community-blogger/api/protos/user"
	"community-blogger/internal/app/user/repositories"
	"community-blogger/internal/pkg/baseerror"
	"community-blogger/internal/pkg/requests"
	"community-blogger/internal/pkg/responses"
	"community-blogger/internal/pkg/utils/middlewareutil"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

var (
	// NotFoundUserErr 用户不存在
	NotFoundUserErr = baseerror.NewBaseError("用户不存在")
	// UserNameOrPasswordErr 用户不存在或者密码错误
	UserNameOrPasswordErr = baseerror.NewBaseError("用户不存在或者密码错误")
	// AccessTokenErr 生成签名错误
	AccessTokenErr = baseerror.NewBaseError("生成签名错误")
	// CreateMemberErr 注册失败
	CreateMemberErr = baseerror.NewBaseError("注册失败")
	// ExistsUserErr 用户已存在，无法注册
	ExistsUserErr = baseerror.NewBaseError("用户已存在，无法注册")
)

// UserService 定义user service
type UserService interface {
	FindByID(ctx context.Context, req *requests.User) (*responses.UserInfo, error)
	FindByToken(ctx context.Context, req *requests.UserToken) (*responses.UserInfo, error)
	Register(req *requests.RegisterRequest) (*responses.RegisterResponse, error)
	Login(req *requests.LoginRequest) (*responses.LoginResponse, error)
}

// DefaultUserService 默认service所拥有的对象
type DefaultUserService struct {
	logger     *zap.Logger
	v          *viper.Viper
	Repository repositories.UserRepository
	userclient userproto.UserClient
}

// NewUserService 初始化UserService
func NewUserService(
	logger *zap.Logger,
	v *viper.Viper,
	repository repositories.UserRepository,
	userclient userproto.UserClient) UserService {
	return &DefaultUserService{
		logger:     logger.With(zap.String("type", "DefaultUserService")),
		v:          v,
		Repository: repository,
		userclient: userclient,
	}
}

// FindByID 调用rpc服务通过ID获取用户信息
func (s *DefaultUserService) FindByID(ctx context.Context, req *requests.User) (*responses.UserInfo, error) {
	rpcReq := userproto.FindByIDRequest{
		Id: int32(req.ID),
	}
	data, err := s.userclient.FindById(ctx, &rpcReq)
	if err != nil {
		s.logger.Error("通过用户ID获取用户调用失败", zap.Error(err))
	}
	res := &responses.UserInfo{
		ID:       int(data.Id),
		Username: data.Username,
		Password: data.Password,
		Token:    data.Token,
	}

	return res, nil
}

// FindByToken 调用rpc服务通过Token获取用户信息
func (s *DefaultUserService) FindByToken(ctx context.Context, req *requests.UserToken) (*responses.UserInfo, error) {
	var res *responses.UserInfo
	rpcReq := userproto.FindByTokenRequest{
		Token: req.Token,
	}
	data, err := s.userclient.FindByToken(ctx, &rpcReq)
	if err != nil {
		s.logger.Error("通过Token获取用户调用失败", zap.Error(err))

	}
	res = &responses.UserInfo{
		ID: int(data.Id),
	}
	return res, nil
}

// Register 用户注册
func (s *DefaultUserService) Register(req *requests.RegisterRequest) (*responses.RegisterResponse, error) {
	mem, _ := s.Repository.FindByUserName(req.Username)
	if mem != nil {
		return nil, ExistsUserErr
	}
	return s.Repository.Register(req)
}

// Login 登录
func (s *DefaultUserService) Login(req *requests.LoginRequest) (*responses.LoginResponse, error) {
	user, err := s.Repository.FindByUserName(req.Username)
	if err != nil {
		return nil, NotFoundUserErr
	}
	if user.Password != fmt.Sprintf("%x", md5.Sum([]byte(req.Password))) {
		return nil, UserNameOrPasswordErr
	}

	expired := time.Now().Add(7 * 24 * time.Hour).Unix()
	accessToken, err := middlewareutil.CreateAccessToken(req.Username, expired)
	if err != nil {
		return nil, AccessTokenErr
	}
	return &responses.LoginResponse{
		Token:       user.Token,
		AccessToken: accessToken,
		ExpireAt:    expired,
		TimeStamp:   time.Now().Unix(),
	}, nil

}
