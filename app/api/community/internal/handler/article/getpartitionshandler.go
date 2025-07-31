package article

import (
	"net/http"

	"github.com/krace-tx/emo_trash/app/api/community/internal/logic/article"
	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 查询文章分区
func GetPartitionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := article.NewGetPartitionsLogic(r.Context(), svcCtx)
		resp, err := l.GetPartitions()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
