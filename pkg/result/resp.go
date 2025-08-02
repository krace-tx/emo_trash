package result

import (
	"github.com/krace-tx/emo_trash/pkg/err"
	"net/http"
)

type Response struct {
	Code    uint32                 `json:"code"`
	Message string                 `json:"message"`
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}

func Ok() *Response {
	return &Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "请求成功",
		Data:    make(map[string]interface{}),
	}
}

func Err() *Response {
	return &Response{
		Code:    http.StatusBadRequest,
		Success: false,
		Message: "请求失败",
		Data:    make(map[string]interface{}),
	}
}

func Errx(err *err.Err) *Response {
	return &Response{
		Code:    err.Code,
		Success: false,
		Message: err.Message,
		Data:    make(map[string]interface{}),
	}
}

func (r *Response) SetMsg(msg string) *Response {
	if r == nil {
		return nil // 防止对 nil 进行操作
	}
	r.Message = msg
	return r
}

func (r *Response) SetCode(code uint32) *Response {
	if r == nil {
		return nil
	}
	r.Code = code
	return r
}

func (r *Response) SetData(key string, value interface{}) *Response {
	if r == nil {
		return nil
	}
	r.Data[key] = value
	return r
}
