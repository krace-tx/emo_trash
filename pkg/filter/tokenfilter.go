package filter

import (
	"context"
	"net/http"
	"zjyt-cloud/pkg/consts"
	"zjyt-cloud/pkg/errorx"
	"zjyt-cloud/pkg/utils"
)

// 该函数从请求头中提取 JWT 令牌，并验证其有效性。
// 如果令牌缺失或无效，则返回 401 未授权错误。
func TokenFilter(r *http.Request, secret string) (*http.Request, error) {
	token := r.Header.Get(consts.Authorize)
	if token == "" {
		// 如果令牌为空，返回未授权响应
		return r, errorx.AuthTokenNotNull
	}

	if token[:7] == "Bearer " {
		token = token[7:]
	}

	// 解析 JWT 令牌，验证其有效性
	claims, err := utils.ParseJwtToken(token, secret)
	if err != nil {
		// 如果令牌验证失败，返回未授权响应
		return r, errorx.AuthTokenFail
	}

	// 将用户 ID 存入请求上下文，并返回更新后的请求
	return r.WithContext(context.WithValue(r.Context(), consts.UserId, claims[consts.UserId])), nil
}
