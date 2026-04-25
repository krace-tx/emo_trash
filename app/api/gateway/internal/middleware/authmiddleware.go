package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"
	consts "github.com/krace-tx/emo_trash/pkg/constant"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type AuthMiddleware struct {
	authRpc auth.Auth
}

func NewAuthMiddleware(authRpc auth.Auth) *AuthMiddleware {
	return &AuthMiddleware{
		authRpc: authRpc,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. 获取 Token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpx.WriteJson(w, http.StatusUnauthorized, types.Unauthorized())
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			httpx.WriteJson(w, http.StatusUnauthorized, types.Unauthorized())
			return
		}
		token := parts[1]

		// 2. 调用 SSO RPC 验证 Token
		resp, err := m.authRpc.VerifyToken(r.Context(), &auth.VerifyReq{
			Token: token,
		})
		if err != nil {
			httpx.WriteJson(w, http.StatusUnauthorized, types.Error(err))
			return
		}

		// 3. 将 user_id 注入 context
		ctx := context.WithValue(r.Context(), consts.UserId, resp.UserId)
		next(w, r.WithContext(ctx))
	}
}
