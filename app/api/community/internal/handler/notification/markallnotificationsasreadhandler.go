package notification

import (
	"net/http"

	"github.com/krace-tx/emo_trash/app/api/community/internal/logic/notification"
	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 标记所有通知为已读
func MarkAllNotificationsAsReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notification.NewMarkAllNotificationsAsReadLogic(r.Context(), svcCtx)
		resp, err := l.MarkAllNotificationsAsRead()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
