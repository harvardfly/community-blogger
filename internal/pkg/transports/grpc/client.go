package grpc

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

////etcd解析器
//type etcdResolver struct {
//	etcdAddr   string
//	clientConn resolver.ClientConn
//}

// ClientOptions grpc client option
type ClientOptions struct {
	Wait            time.Duration
	Tag             string
	GrpcDialOptions []grpc.DialOption
}

// NewClientOptions new grpc client option
func NewClientOptions(v *viper.Viper, logger *zap.Logger) (*ClientOptions, error) {
	var (
		err error
		o   = new(ClientOptions)
	)
	if err = v.UnmarshalKey("grpc.client", o); err != nil {
		return nil, err
	}

	logger.Info("load grpc.client options success", zap.Any("grpc.client options", o))

	return o, nil
}

// ClientOptional grpc client optional
type ClientOptional func(o *ClientOptions)

// WithTimeout grpc client time out
func WithTimeout(d time.Duration) ClientOptional {
	return func(o *ClientOptions) {
		o.Wait = d
	}
}

// WithTag grpc client tag
func WithTag(tag string) ClientOptional {
	return func(o *ClientOptions) {
		o.Tag = tag
	}
}

// WithGrpcDialOptions grpc dial option
func WithGrpcDialOptions(options ...grpc.DialOption) ClientOptional {
	return func(o *ClientOptions) {
		o.GrpcDialOptions = append(o.GrpcDialOptions, options...)
	}
}

// Client grpc client server
type Client struct {
	o *ClientOptions
}

// NewClient new grpc client server
func NewClient(o *ClientOptions) (*Client, error) {
	return &Client{
		o: o,
	}, nil
}

// Dial grpc client dial
func (c *Client) Dial(service string, options ...ClientOptional) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	o := &ClientOptions{
		Wait:            c.o.Wait,
		Tag:             c.o.Tag,
		GrpcDialOptions: c.o.GrpcDialOptions,
	}
	for _, option := range options {
		option(o)
	}

	conn, err := grpc.DialContext(ctx, service, o.GrpcDialOptions...)
	if err != nil {
		return nil, errors.Wrap(err, "grpc dial error")
	}

	return conn, nil
}

////初始化一个etcd解析器
//func NewResolver(etcdAddr string) resolver.Builder {
//	return &etcdResolver{etcdAddr: etcdAddr}
//}
//
//func (r *etcdResolver) Scheme() string {
//	return schema
//}
//
////watch有变化以后会调用
//func (r *etcdResolver) ResolveNow(rn resolver.ResolveNowOptions) {
//	log.Println("ResolveNow")
//	fmt.Println(rn)
//}
//
////解析器关闭时调用
//func (r *etcdResolver) Close() {
//	log.Println("Close")
//}
//
////构建解析器 grpc.Dial()同步调用
//func (r *etcdResolver) Build(target resolver.Target, clientConn resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
//	var err error
//
//	//构建etcd client
//	if cli == nil {
//		cli, err = clientv3.New(clientv3.Config{
//			Endpoints:   strings.Split(r.etcdAddr, ";"),
//			DialTimeout: 15 * time.Second,
//		})
//		if err != nil {
//			fmt.Printf("连接etcd失败：%s\n", err)
//			return nil, err
//		}
//	}
//
//	r.clientConn = clientConn
//
//	go r.watch("/" + target.Scheme + "/" + target.Endpoint + "/")
//
//	return r, nil
//}
//
////监听etcd中某个key前缀的服务地址列表的变化
//func (r *etcdResolver) watch(keyPrefix string) {
//	//初始化服务地址列表
//	var addrList []resolver.Address
//
//	resp, err := cli.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
//	if err != nil {
//		fmt.Println("获取服务地址列表失败：", err)
//	} else {
//		for i := range resp.Kvs {
//			addrList = append(addrList, resolver.Address{Addr: strings.TrimPrefix(string(resp.Kvs[i].Key), keyPrefix)})
//		}
//	}
//
//	r.clientConn.NewAddress(addrList)
//
//	//监听服务地址列表的变化
//	rch := cli.Watch(context.Background(), keyPrefix, clientv3.WithPrefix())
//	for n := range rch {
//		for _, ev := range n.Events {
//			addr := strings.TrimPrefix(string(ev.Kv.Key), keyPrefix)
//			switch ev.Type {
//			case mvccpb.PUT:
//				if !exists(addrList, addr) {
//					addrList = append(addrList, resolver.Address{Addr: addr})
//					r.clientConn.NewAddress(addrList)
//				}
//			case mvccpb.DELETE:
//				if s, ok := remove(addrList, addr); ok {
//					addrList = s
//					r.clientConn.NewAddress(addrList)
//				}
//			}
//		}
//	}
//}
//
//func exists(l []resolver.Address, addr string) bool {
//	for i := range l {
//		if l[i].Addr == addr {
//			return true
//		}
//	}
//	return false
//}
//
//func remove(s []resolver.Address, addr string) ([]resolver.Address, bool) {
//	for i := range s {
//		if s[i].Addr == addr {
//			s[i] = s[len(s)-1]
//			return s[:len(s)-1], true
//		}
//	}
//	return nil, false
//}
