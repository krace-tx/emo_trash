package sso

import (
	"context"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendSmsCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendSmsCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendSmsCodeLogic {
	return &SendSmsCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendSmsCodeLogic) SendSmsCode(req *types.SendSmsCodeReq) (resp *types.CommonResp, err error) {

	data, err := l.svcCtx.Auth.SendSmsCode(l.ctx, &auth.SendSmsCodeReq{
		Mobile: req.Mobile,
		Scene:  req.Scene,
	})
	if err != nil {
		l.Logger.Errorf("SendSmsCode failed, err: %v", err)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
