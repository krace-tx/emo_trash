package sso

import (
	"context"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.Auth.LoginByPassword(l.ctx, &auth.LoginByPasswordReq{
		Account:    req.Account,
		Password:   req.Password,
		DeviceType: req.DeviceType,
		DeviceId:   req.DeviceId,
		LoginIp:    req.LoginIp,
	})
	if err != nil {
		l.Logger.Errorf("Login failed, err: %v", err)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
