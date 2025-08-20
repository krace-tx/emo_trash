package interceptor

import (
	"context"
	"fmt"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"regexp"
	"strings"
)

type Validator interface {
	Validate() error
	ValidateAll() error
}

type MultiValidationError interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	Error() string
}

// UnaryValidatorInterceptor 自定义验证拦截器
func UnaryValidatorInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if validator, ok := req.(Validator); ok {
		if err := validator.Validate(); err != nil {
			if err := translateValidationError(req, err); err != nil {
				return nil, err
			}
			return nil, errx.ErrSystemArgInvalid
		}
	}
	return handler(ctx, req)
}

// translateValidationError 将验证错误转换为自定义描述
func translateValidationError(req any, err error) error {
	if err == nil {
		return nil
	}

	// 处理多字段验证错误
	if multiErr, ok := err.(MultiValidationError); ok {
		return translateMultiValidationError(req, multiErr)
	}

	return nil
}

// translateMultiValidationError 处理多字段验证错误
func translateMultiValidationError(req any, multiErr MultiValidationError) error {
	var errorMsgs []string

	// 递归处理嵌套错误
	if cause := multiErr.Cause(); cause != nil {
		if nestedMultiErr, ok := cause.(MultiValidationError); ok {
			return translateMultiValidationError(req, nestedMultiErr)
		}
		return translateValidationError(req, cause)
	}

	// 获取字段和原因并翻译
	field := multiErr.Field()

	translatedReason := translateReason(req, field)
	if translatedReason != nil {
		e := translatedReason.Error()

		errorMsgs = append(errorMsgs, e)
	}

	return fmt.Errorf(strings.Join(errorMsgs, "; "))
}

// translateReason 根据错误原因和字段名生成中文错误描述
func translateReason(req any, field string) error {
	pb, ok := req.(proto.Message)
	if !ok {
		return nil
	}
	desc := pb.ProtoReflect().Descriptor()
	if desc == nil {
		return nil
	}

	fields := desc.Fields()

	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		if matchField(fd.TextName(), field) {
			re := regexp.MustCompile(`50000:"([^"]*)"`)
			optionsStr := fmt.Sprintf("%v", fd.Options())
			matches := re.FindStringSubmatch(optionsStr)
			if len(matches) > 1 && matches[1] != "" {
				return fmt.Errorf("%s: %s", field, matches[1])
			}
		}
	}

	return nil
}

func matchField(fd1, fd2 string) bool {
	// 统一转换函数：移除下划线并转为小写
	normalize := func(s string) string {
		// 移除所有下划线
		s = strings.ReplaceAll(s, "_", "")
		// 转换为小写
		return strings.ToLower(s)
	}

	// 比较标准化后的字段名
	return normalize(fd1) == normalize(fd2)
}
