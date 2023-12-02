package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

func UnaryTimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if timeout <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		ctx, cancelFunc := context.WithTimeout(ctx, timeout)
		defer cancelFunc()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
