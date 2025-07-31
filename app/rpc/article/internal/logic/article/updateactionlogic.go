package articlelogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/article/article"
	"github.com/krace-tx/emo_trash/app/rpc/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateActionLogic {
	return &UpdateActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新活动状态
func (l *UpdateActionLogic) UpdateAction(in *article.UpdateActionRequest) (*article.UpdateActionResponse, error) {
	// todo: add your logic here and delete this line

	return &article.UpdateActionResponse{}, nil
}
