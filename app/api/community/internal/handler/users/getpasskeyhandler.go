package users

import (
	"net/http"

	"github.com/krace-tx/emo_trash/app/api/community/internal/logic/users"
	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取一次性通行证
func GetPassKeyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := users.NewGetPassKeyLogic(r.Context(), svcCtx)
		resp, err := l.GetPassKey()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
