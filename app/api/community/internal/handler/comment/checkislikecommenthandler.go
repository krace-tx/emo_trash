package comment

import (
	"net/http"

	"github.com/krace-tx/emo_trash/app/api/community/internal/logic/comment"
	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 检查用户是否对评论进行点赞
func CheckIsLikeCommentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckIsLikeCommentReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := comment.NewCheckIsLikeCommentLogic(r.Context(), svcCtx)
		resp, err := l.CheckIsLikeComment(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
