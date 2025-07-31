package users

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPassKeyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取一次性通行证
func NewGetPassKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPassKeyLogic {
	return &GetPassKeyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPassKeyLogic) GetPassKey() (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
