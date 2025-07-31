package chatlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/chat/chat"
	"github.com/krace-tx/emo_trash/app/rpc/chat/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatHistoryLogic {
	return &GetChatHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询聊天记录
func (l *GetChatHistoryLogic) GetChatHistory(in *chat.GetChatHistoryRequest) (*chat.GetChatHistoryResponse, error) {
	// todo: add your logic here and delete this line

	return &chat.GetChatHistoryResponse{}, nil
}
