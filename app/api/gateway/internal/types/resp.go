package types

import errx "github.com/krace-tx/emo_trash/pkg/err"

type Builder struct {
	resp *CommonResp
}

func New() *Builder {
	return &Builder{resp: &CommonResp{}}
}

func (b *Builder) Code(code uint32) *Builder {
	b.resp.Code = code
	return b
}

func (b *Builder) Success(success bool) *Builder {
	b.resp.Success = success
	return b
}

func (b *Builder) Data(data interface{}) *Builder {
	b.resp.Data = data
	return b
}

func (b *Builder) Message(msg string) *Builder {
	b.resp.Message = msg
	return b
}

func (b *Builder) Build() *CommonResp {
	return b.resp
}

// Success 标准成功响应（带数据)
func Success(data interface{}) *CommonResp {
	return New().
		Code(200).
		Success(true).
		Data(data).
		Message("success").
		Build()
}

// SuccessWithMsg 带自定义消息成功响应
func SuccessWithMsg(data interface{}, msg string) *CommonResp {
	return New().
		Code(200).
		Success(true).
		Data(data).
		Message(msg).
		Build()
}

// SuccessEmpty 无数据成功响应
func SuccessEmpty() *CommonResp {
	return New().
		Code(200).
		Success(true).
		Message("success").
		Build()
}

// Error 通用错误响应
func Error(err error) *CommonResp {
	parseError := errx.ParseError(err)

	return New().
		Code(parseError.Code).
		Success(false).
		Message(parseError.Message).
		Build()
}

// ParamError 参数错误响应（400）
func ParamError(msg string) *CommonResp {
	if msg == "" {
		msg = "invalid parameters"
	}
	return New().
		Code(400).
		Success(false).
		Message(msg).
		Build()
}

// Unauthorized 未授权响应（401）
func Unauthorized() *CommonResp {
	return New().
		Code(401).
		Success(false).
		Message("authentication required").
		Build()
}

// Forbidden 权限拒绝响应（403）
func Forbidden() *CommonResp {
	return New().
		Code(403).
		Success(false).
		Message("permission denied").
		Build()
}

// NotFound 资源不存在响应（404）
func NotFound(resource string) *CommonResp {
	msg := "resource not found"
	if resource != "" {
		msg = resource + " not found"
	}
	return New().
		Code(404).
		Success(false).
		Message(msg).
		Build()
}

// ServerError 服务器内部错误响应（500）
func ServerError() *CommonResp {
	return New().
		Code(500).
		Success(false).
		Message("internal server error").
		Build()
}
