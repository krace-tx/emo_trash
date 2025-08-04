package chatlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/chat/chat"
	"github.com/krace-tx/emo_trash/app/rpc/chat/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAiChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiChatLogic {
	return &AiChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AI 聊天接口（新增）
func (l *AiChatLogic) AiChat(in *chat.AIChatRequest) (*chat.AIChatResponse, error) {
	// todo: add your logic here and delete this line

	return &chat.AIChatResponse{}, nil
}
