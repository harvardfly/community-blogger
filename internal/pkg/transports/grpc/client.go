package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ClientOptions grpc client option
type ClientOptions struct {
	Wait            time.Duration
	Tag             string
	GrpcDialOptions []grpc.DialOption
}

func init() {
	DefaultDialer.CAFile = "../internal/pkg/transports/tls/client/ca.pem"
	DefaultDialer.CertFile = "../internal/pkg/transports/tls/client/client.pem"
	DefaultDialer.KeyFile = "../internal/pkg/transports/tls/client/client.key"
}

// GrpcDialer .
type Dialer struct {
	CertFile string
	KeyFile  string
	CAFile   string
}

var DefaultDialer = Dialer{}

// TransportCredentials 获取传输TSL证书
func (d Dialer) TransportCredentials() (grpc.DialOption, error) {
	cert, err := tls.LoadX509KeyPair(d.CertFile, d.KeyFile)
	if err != nil {
		return nil, err
	}
	ca, err := ioutil.ReadFile(d.CAFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, errors.New("failed to append certs from pem")
	}

	tc := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	})
	return grpc.WithTransportCredentials(tc), nil
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
	// TLS安全验证
	credential, _ := DefaultDialer.TransportCredentials()
	if credential != nil {
		o.GrpcDialOptions = append(o.GrpcDialOptions, credential)
	} else {
		// 默认不加安全验证
		options = append(options, WithGrpcDialOptions(grpc.WithInsecure()))
	}
	for _, option := range options {
		option(o)
	}
	fmt.Println(service)
	fmt.Println("00000000000000000000")
	fmt.Println(o.GrpcDialOptions)
	conn, err := grpc.DialContext(ctx, service, o.GrpcDialOptions...)
	if err != nil {
		fmt.Println("11111111111111111111111111")
		return nil, errors.Wrap(err, "grpc dial error")
	}
	fmt.Println("2222222222222222")
	fmt.Println(conn)
	return conn, nil
}
