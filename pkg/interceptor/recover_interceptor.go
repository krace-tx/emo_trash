package interceptor

import (
	"context"
	"runtime/debug"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryRecoverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	defer func() {
		if r := recover(); r != nil {
			logx.WithContext(ctx).Errorf("panic recovered: %v", r)
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "panic: %v", r)
		}
	}()

	return handler(ctx, req)
}
