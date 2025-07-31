package notification

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkAllNotificationsAsReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标记所有通知为已读
func NewMarkAllNotificationsAsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkAllNotificationsAsReadLogic {
	return &MarkAllNotificationsAsReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkAllNotificationsAsReadLogic) MarkAllNotificationsAsRead() (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
