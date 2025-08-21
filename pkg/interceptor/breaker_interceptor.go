package interceptor

import (
	"context"
	"github.com/zeromicro/go-zero/core/breaker"
	"google.golang.org/grpc"
)

func BreakerInterceptor() grpc.UnaryClientInterceptor {
	brk := breaker.NewBreaker()
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return brk.DoWithAcceptable(func() error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}, func(err error) bool {
			// 定义哪些错误不计入失败率
			return err == nil
		})
	}
}
