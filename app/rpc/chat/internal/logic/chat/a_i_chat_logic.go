package chatlogic

import (
	"context"
	"github.com/krace-tx/emo_trash/app/rpc/chat/chat"
	"github.com/krace-tx/emo_trash/app/rpc/chat/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AIChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAIChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AIChatLogic {
	return &AIChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AI 聊天接口（新增）
func (l *AIChatLogic) AIChat(in *chat.AIChatRequest) (*chat.AIChatResponse, error) {

	return &chat.AIChatResponse{
		ReplyText:      "",
		ConversationId: "",
		Sentiment:      nil,
		Timestamp:      0,
		RequestId:      "",
		DebugInfo:      nil,
	}, nil
}
