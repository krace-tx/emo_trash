package notification

import (
	"net/http"

	"github.com/krace-tx/emo_trash/app/api/community/internal/logic/notification"
	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取未读通知数量
func GetUnreadNotificationCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := notification.NewGetUnreadNotificationCountLogic(r.Context(), svcCtx)
		resp, err := l.GetUnreadNotificationCount()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
