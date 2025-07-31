package notificationlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/notification/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/notification/notification"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAllNotificationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAllNotificationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAllNotificationsLogic {
	return &DeleteAllNotificationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除所有通知
func (l *DeleteAllNotificationsLogic) DeleteAllNotifications(in *notification.DeleteAllNotificationsRequest) (*notification.DeleteAllNotificationsResponse, error) {
	// todo: add your logic here and delete this line

	return &notification.DeleteAllNotificationsResponse{}, nil
}
