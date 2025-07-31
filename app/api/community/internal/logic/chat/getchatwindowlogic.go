package chat

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatWindowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询指定对话窗口的聊天记录
func NewGetChatWindowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatWindowLogic {
	return &GetChatWindowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChatWindowLogic) GetChatWindow(req *types.GetChatWindowReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
