package chatlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/chat/chat"
	"github.com/krace-tx/emo_trash/app/rpc/chat/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOnlineStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOnlineStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOnlineStatusLogic {
	return &UpdateOnlineStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新在线状态
func (l *UpdateOnlineStatusLogic) UpdateOnlineStatus(in *chat.UpdateOnlineStatusRequest) (*chat.UpdateOnlineStatusResponse, error) {
	// todo: add your logic here and delete this line

	return &chat.UpdateOnlineStatusResponse{}, nil
}
