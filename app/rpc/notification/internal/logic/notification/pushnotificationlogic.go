package notificationlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/notification/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/notification/notification"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushNotificationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushNotificationLogic {
	return &PushNotificationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 推送通知
func (l *PushNotificationLogic) PushNotification(in *notification.PushNotificationRequest) (*notification.PushNotificationResponse, error) {
	// todo: add your logic here and delete this line

	return &notification.PushNotificationResponse{}, nil
}
