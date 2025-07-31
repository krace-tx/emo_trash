package notification

import (
	"net/http"

	"github.com/krace-tx/emo_trash/app/api/community/internal/logic/notification"
	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除所有通知
func DeleteAllNotificationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notification.NewDeleteAllNotificationsLogic(r.Context(), svcCtx)
		resp, err := l.DeleteAllNotifications()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
