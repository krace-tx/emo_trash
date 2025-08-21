package sso

import (
	"context"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateQrcodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateQrcodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateQrcodeLogic {
	return &GenerateQrcodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateQrcodeLogic) GenerateQrcode(req *types.QrcodeReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.Auth.GenerateQrcode(l.ctx, &auth.QrcodeReq{
		DeviceId: req.DeviceId,
		LoginIp:  req.LoginIp,
	})
	if err != nil {
		l.Logger.Errorf("GenerateQrcode failed, err: %v", err)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
