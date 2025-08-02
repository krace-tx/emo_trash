package notification

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushNotificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 推送通知
func NewPushNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushNotificationLogic {
	return &PushNotificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PushNotificationLogic) PushNotification(req *types.PushNotificationReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
