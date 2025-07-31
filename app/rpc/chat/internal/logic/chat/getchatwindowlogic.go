package chatlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/chat/chat"
	"github.com/krace-tx/emo_trash/app/rpc/chat/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatWindowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatWindowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatWindowLogic {
	return &GetChatWindowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询指定对话窗口的聊天记录
func (l *GetChatWindowLogic) GetChatWindow(in *chat.GetChatWindowRequest) (*chat.GetChatWindowResponse, error) {
	// todo: add your logic here and delete this line

	return &chat.GetChatWindowResponse{}, nil
}
