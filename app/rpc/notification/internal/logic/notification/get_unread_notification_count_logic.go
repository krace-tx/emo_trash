package notificationlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/notification/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/notification/notification"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUnreadNotificationCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUnreadNotificationCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnreadNotificationCountLogic {
	return &GetUnreadNotificationCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取未读通知数量
func (l *GetUnreadNotificationCountLogic) GetUnreadNotificationCount(in *notification.GetUnreadNotificationCountRequest) (*notification.GetUnreadNotificationCountResponse, error) {
	// todo: add your logic here and delete this line

	return &notification.GetUnreadNotificationCountResponse{}, nil
}
