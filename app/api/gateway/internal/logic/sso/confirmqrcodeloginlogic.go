package sso

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmQrcodeLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfirmQrcodeLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmQrcodeLoginLogic {
	return &ConfirmQrcodeLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfirmQrcodeLoginLogic) ConfirmQrcodeLogin(req *types.QrcodeConfirmReq) (resp *types.CommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
