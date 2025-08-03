package chat

import (
	"github.com/krace-tx/emo_trash/pkg/result"
	"net/http"

	"github.com/krace-tx/emo_trash/app/api/community/internal/logic/chat"
	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// AI 聊天机器人
func ChatBotHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendMessageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := chat.NewChatBotLogic(r.Context(), svcCtx)
		resp, err := l.ChatBot(&req)
		if err != nil {
			result.ErrorCtx(r.Context(), w, err)
		} else {
			result.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
