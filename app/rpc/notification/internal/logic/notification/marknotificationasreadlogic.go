package notificationlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/notification/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/notification/notification"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkNotificationAsReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkNotificationAsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkNotificationAsReadLogic {
	return &MarkNotificationAsReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 标记通知为已读
func (l *MarkNotificationAsReadLogic) MarkNotificationAsRead(in *notification.MarkNotificationAsReadRequest) (*notification.MarkNotificationAsReadResponse, error) {
	// todo: add your logic here and delete this line

	return &notification.MarkNotificationAsReadResponse{}, nil
}
