package sso

import (
	"context"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendEmailCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendEmailCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendEmailCodeLogic {
	return &SendEmailCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendEmailCodeLogic) SendEmailCode(req *types.SendEmailCodeReq) (resp *types.CommonResp, err error) {

	data, err := l.svcCtx.Auth.SendEmailCode(l.ctx, &auth.SendEmailCodeReq{
		Email: req.Email,
		Scene: req.Scene,
	})
	if err != nil {
		l.Logger.Errorf("SendEmailCode failed, err: %v", err)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
