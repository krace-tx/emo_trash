package result

import (
	"context"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

type Response struct {
	Code    uint32      `json:"code"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func ErrorCtx(context context.Context, w http.ResponseWriter, err error) {
	httpx.ErrorCtx(context, w, err)
}

func OkJsonCtx(context context.Context, w http.ResponseWriter, resp interface{}) {
	response := &Response{
		Code:    http.StatusOK,
		Message: "success",
		Success: true,
		Data:    resp,
	}
	httpx.OkJsonCtx(context, w, response)
}
