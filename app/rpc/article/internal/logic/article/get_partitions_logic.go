package articlelogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/article/article"
	"github.com/krace-tx/emo_trash/app/rpc/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPartitionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPartitionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPartitionsLogic {
	return &GetPartitionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询分区
func (l *GetPartitionsLogic) GetPartitions(in *article.GetPartitionsRequest) (*article.GetPartitionsResponse, error) {
	// todo: add your logic here and delete this line

	return &article.GetPartitionsResponse{}, nil
}
