package grpcclients

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	grpc2 "google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"log"
	userproto "community-blogger/api/protos/user"
	"community-blogger/internal/pkg/etcdservice"
	"community-blogger/internal/pkg/transports/grpc"
)

// NewUserClient 初始化 user rpc client
func NewUserClient(client *grpc.Client, v *viper.Viper) (userproto.UserClient, error) {
	o := new(clientTarget)
	if err := v.UnmarshalKey("client.target", o); err != nil {
		return nil, errors.Wrap(err, "获取client.target配置失败")
	}
	//注册etcd解析器
	r := etcdservice.NewResolver(o.EtcdAddr)
	resolver.Register(r)

	// 客户端连接服务器
	conn, err := grpc2.Dial(r.Scheme()+"://"+o.Caller+"/"+o.Callee, grpc2.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)), grpc2.WithInsecure())

	if err != nil {
		log.Println("连接服务器失败", err)
		return nil, errors.Wrap(err, "notify client dial error")
	}

	return userproto.NewUserClient(conn), nil
}
