package sso

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindMobileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBindMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindMobileLogic {
	return &BindMobileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindMobileLogic) BindMobile(req *types.BindMobileReq) (resp *types.CommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
