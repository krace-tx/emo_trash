package interceptor

import (
	"context"
	"fmt"
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
			// 需要捕获的常见panic信号类型
			// 1. 空指针解引用 (nil pointer dereference)
			// 2. 数组/切片索引越界 (index out of range)
			// 3. 类型断言失败 (type assertion failed)
			// 4. 关闭的通道发送数据 (send on closed channel)
			// 5. 除以零 (division by zero)
			// 6. 无效的内存访问 (invalid memory address)
			// 7. 自定义panic（如业务代码中主动调用panic()）

			// 返回标准gRPC错误
			err = status.Errorf(codes.Internal, "system error")

			// 分类处理不同类型的panic值
			var panicType string
			switch v := r.(type) {
			case error:
				panicType = "error"
				logx.WithContext(ctx).Errorf("panic recovered [type:%s, msg:%v]", panicType, v.Error())
			case string:
				panicType = "string"
				logx.WithContext(ctx).Errorf("panic recovered [type:%s, msg:%s]", panicType, v)
			default:
				panicType = fmt.Sprintf("unknown(%T)", v)
				logx.WithContext(ctx).Errorf("panic recovered [type:%s, value:%v]", panicType, v)
			}

			// 打印调用栈信息
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "panic: %v", r)
		}
	}()

	return handler(ctx, req)
}
