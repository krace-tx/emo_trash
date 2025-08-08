package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryAuthInterceptor 自定义鉴权拦截器
// 自定义 auth 捕获拦截器、框架中如果 auth 设置为 true，自定义拦截器会失效
func UnaryAuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	// 1. 从 metadata 获取 token
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	// 2. 校验 token（示例：检查是否包含 "my-secret-token"）
	tokens := md.Get("authorization")
	if len(tokens) == 0 || tokens[0] != "my-secret-token" {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// 3. 继续处理请求
	return handler(ctx, req)
}

// 自定义 Stream 拦截器
func StreamAuthInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// 1. 从上下文获取 metadata（如 token）
	ctx := ss.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "missing metadata")
	}

	// 2. 校验 token
	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	// 3. 继续处理流
	return handler(srv, ss)
}
