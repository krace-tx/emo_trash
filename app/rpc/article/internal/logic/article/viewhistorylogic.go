package articlelogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/article/article"
	"github.com/krace-tx/emo_trash/app/rpc/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ViewHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewViewHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ViewHistoryLogic {
	return &ViewHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询历史记录
func (l *ViewHistoryLogic) ViewHistory(in *article.ViewHistoryRequest) (*article.ViewHistoryResponse, error) {
	// todo: add your logic here and delete this line

	return &article.ViewHistoryResponse{}, nil
}
