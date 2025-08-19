//package interceptor
//
//import (
//	"context"
//	"fmt"
//	"github.com/go-playground/validator/v10"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/status"
//	"strings"
//
//	"google.golang.org/protobuf/proto"
//)
//
//type Validator interface {
//	Validator() error
//}
//
//// transplateToChinese 默认的中文翻译
//func (cv *Validator) translateToChinese(err validator.FieldError) string {
//	fieldName := cv.getChineseFieldName(err.Field())
//
//	switch err.Tag() {
//	case "required":
//		return fmt.Sprintf("%s是必填字段", fieldName)
//	case "min":
//		return fmt.Sprintf("%s长度不能少于%s个字符", fieldName, err.Param())
//	case "max":
//		return fmt.Sprintf("%s长度不能超过%s个字符", fieldName, err.Param())
//	case "len":
//		return fmt.Sprintf("%s长度必须为%s个字符", fieldName, err.Param())
//	case "email":
//		return fmt.Sprintf("%s必须是有效的邮箱格式", fieldName)
//	default:
//		return fmt.Sprintf("%s验证失败", fieldName)
//	}
//}
//
//// getChineseFieldName 获取字段的中文名
//func (cv *Validator) getChineseFieldName(field string) string {
//	mappings := map[string]string{
//		"name":     "姓名",
//		"email":    "邮箱",
//		"password": "密码",
//		"mobile":   "手机号",
//		"phone":    "手机号",
//		"age":      "年龄",
//	}
//
//	if chs, exists := mappings[strings.ToLower(field)]; exists {
//		return chs
//	}
//	return field
//}
//
//// UnaryValidatorInterceptor 自定义验证拦截器
//func UnaryValidatorInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
//	return func {
//		// 检查是否是 proto.Message
//		if msg, ok := req.(proto.Message); ok {
//			if err := validator.ValidateWithCustomMessages(msg); err != nil {
//				return nil, status.Errorf(codes.InvalidArgument, "参数验证失败: %s", err.Error())
//			}
//		}
//
//		return handler(ctx, req)
//	}
//}
