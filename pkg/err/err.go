package errx

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"regexp"
	"strconv"
	"strings"
)

// errorx 结构体定义，包含错误码和错误信息
type Err struct {
	Code    uint32 // 错误码
	Message string // 错误信息
}

// 实现 error 接口
func (e *Err) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// ParseError 解析错误字符串，返回 Code 和 Message
func ParseError(err error) *Err {
	// 正则表达式匹配错误字符串
	errStr := err.Error()
	re := regexp.MustCompile(`Code: (\d+), Message: (.+)`)
	matches := re.FindStringSubmatch(errStr)

	if len(matches) != 3 {
		return &Err{Error.Code, fmt.Sprintf("无法解析错误字符串: %s", errStr)}
	}

	code, err := strconv.Atoi(matches[1])
	if err != nil {
		return &Err{Error.Code, fmt.Sprintf("无效的错误代码: %s", matches[1])}
	}

	message := strings.TrimSpace(matches[2])
	return &Err{Code: uint32(code), Message: message}
}

// 创建新的 errorx 实例
func New(code uint32, message string) *Err {
	return &Err{
		Code:    code,
		Message: message,
	}
}

func Errs(message string) *Err {
	logx.Errorf(message)
	return &Err{
		Code:    Error.Code,
		Message: message,
	}
}

func Errf(format string, v ...any) *Err {
	logx.Errorf(format, v)
	return &Err{
		Code:    Error.Code,
		Message: fmt.Sprintf(format, v),
	}
}
