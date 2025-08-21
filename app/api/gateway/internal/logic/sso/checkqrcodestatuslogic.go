package sso

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckQrcodeStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckQrcodeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckQrcodeStatusLogic {
	return &CheckQrcodeStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckQrcodeStatusLogic) CheckQrcodeStatus(req *types.QrcodeStatusReq) (resp *types.CommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
