package notification

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUnreadNotificationCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取未读通知数量
func NewGetUnreadNotificationCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnreadNotificationCountLogic {
	return &GetUnreadNotificationCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUnreadNotificationCountLogic) GetUnreadNotificationCount() (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
