package notificationlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/notification/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/notification/notification"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkAllNotificationsAsReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkAllNotificationsAsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkAllNotificationsAsReadLogic {
	return &MarkAllNotificationsAsReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 标记所有通知为已读
func (l *MarkAllNotificationsAsReadLogic) MarkAllNotificationsAsRead(in *notification.MarkAllNotificationsAsReadRequest) (*notification.MarkAllNotificationsAsReadResponse, error) {
	// todo: add your logic here and delete this line

	return &notification.MarkAllNotificationsAsReadResponse{}, nil
}
