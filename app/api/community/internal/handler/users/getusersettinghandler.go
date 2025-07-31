package users

import (
	"net/http"

	"github.com/krace-tx/emo_trash/app/api/community/internal/logic/users"
	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取用户隐私权限
func GetUserSettingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserSettingReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := users.NewGetUserSettingLogic(r.Context(), svcCtx)
		resp, err := l.GetUserSetting(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
