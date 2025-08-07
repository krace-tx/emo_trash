package sso

import (
	"net/http"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/logic/sso"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UnbindThirdPartyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UnbindThirdPartyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := sso.NewUnbindThirdPartyLogic(r.Context(), svcCtx)
		resp, err := l.UnbindThirdParty(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
