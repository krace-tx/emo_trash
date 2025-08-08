package interceptor

import (
	"context"
	"runtime/debug"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryRecoverInterceptor 自定义恢复拦截器
// 自定义 recover 捕获拦截器、框架中如果 recover 设置为 true，自定义拦截器会失效
func UnaryRecoverInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			logx.WithContext(ctx).Errorf("panic recovered: %v", r)
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "panic: %v", r)
		}
	}()

	return handler(ctx, req)
}
