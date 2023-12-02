package grpcclient

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
	"github.com/duanshanghanqing/rocket/client/grpcclient/interceptors"
)

type GrpcClient struct {
	target                   string                         // address “192.168.0.1:8090”
	timeout                  time.Duration                  // timeout
	unaryClientInterceptors  []grpc.UnaryClientInterceptor  // interceptor
	streamClientInterceptors []grpc.StreamClientInterceptor // stream interceptor
	dialOptions              []grpc.DialOption              // grpc Client Connection Options
	balancerName             string                         // Debt balancing strategy name
	insecure                 bool                           // Do you use unsafe connections
}

type GrpcClientOption func(c *GrpcClient)

func WithGrpcClientOptionTarget(target string) GrpcClientOption {
	return func(c *GrpcClient) {
		c.target = target
	}
}

func WithGrpcClientOptionTimeout(timeout time.Duration) GrpcClientOption {
	return func(c *GrpcClient) {
		c.timeout = timeout
	}
}

func WithGrpcClientOptionUnaryClientInterceptors(unaryClientInterceptors []grpc.UnaryClientInterceptor) GrpcClientOption {
	return func(c *GrpcClient) {
		c.unaryClientInterceptors = unaryClientInterceptors
	}
}

func WithGrpcClientOptionStreamClientInterceptors(streamClientInterceptors []grpc.StreamClientInterceptor) GrpcClientOption {
	return func(c *GrpcClient) {
		c.streamClientInterceptors = streamClientInterceptors
	}
}

func WithGrpcClientOptionDialOption(dialOptions []grpc.DialOption) GrpcClientOption {
	return func(c *GrpcClient) {
		c.dialOptions = dialOptions
	}
}

func WithGrpcClientOptionBalancerName(balancerName string) GrpcClientOption {
	return func(c *GrpcClient) {
		c.balancerName = balancerName
	}
}

func WithGrpcClientOptionInsecure(insecure bool) GrpcClientOption {
	return func(c *GrpcClient) {
		c.insecure = insecure
	}
}

func New(ctx context.Context, opts ...GrpcClientOption) (*grpc.ClientConn, error) {
	c := &GrpcClient{
		timeout:                  5 * time.Second,
		balancerName:             "round_robin", // Set polling algorithm
		insecure:                 true,          // Default use of insecure transmission
		unaryClientInterceptors:  make([]grpc.UnaryClientInterceptor, 0),
		streamClientInterceptors: make([]grpc.StreamClientInterceptor, 0),
		dialOptions:              make([]grpc.DialOption, 0),
	}

	for _, o := range opts {
		o(c)
	}

	var dialOptions = []grpc.DialOption{
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy": "%s"}`, c.balancerName)),
	}

	if c.insecure {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	c.unaryClientInterceptors = append(c.unaryClientInterceptors, interceptors.UnaryTimeoutInterceptor(c.timeout))
	if len(c.unaryClientInterceptors) > 0 {
		dialOptions = append(dialOptions, grpc.WithChainUnaryInterceptor(c.unaryClientInterceptors...))
	}

	if len(c.streamClientInterceptors) > 0 {
		dialOptions = append(dialOptions, grpc.WithChainStreamInterceptor(c.streamClientInterceptors...))
	}

	dialOptions = append(dialOptions, c.dialOptions...)

	return grpc.DialContext(context.Background(), c.target, dialOptions...)
}
