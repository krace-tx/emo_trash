package social

import (
	"net/http"

	"github.com/krace-tx/emo_trash/app/api/community/internal/logic/social"
	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 取消关注
func UnfollowUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UnfollowUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := social.NewUnfollowUserLogic(r.Context(), svcCtx)
		resp, err := l.UnfollowUser(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
