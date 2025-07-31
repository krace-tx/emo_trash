package chat

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询聊天记录
func NewGetChatHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatHistoryLogic {
	return &GetChatHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChatHistoryLogic) GetChatHistory(req *types.GetChatHistoryReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
