package users

import (
	"net/http"

	"github.com/krace-tx/emo_trash/app/api/community/internal/logic/users"
	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 编辑用户隐私权限
func EditUserSettingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EditUserSettingReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := users.NewEditUserSettingLogic(r.Context(), svcCtx)
		resp, err := l.EditUserSetting(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
