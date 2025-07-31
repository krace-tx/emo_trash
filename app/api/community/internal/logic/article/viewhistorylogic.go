package article

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ViewHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询历史记录
func NewViewHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ViewHistoryLogic {
	return &ViewHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ViewHistoryLogic) ViewHistory(req *types.ViewHistoryReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
