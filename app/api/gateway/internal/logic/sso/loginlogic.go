// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

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
	data, err := l.svcCtx.Auth.Login(l.ctx, &auth.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		l.Logger.Errorf("邮箱登录失败: %v, email=%s", err, req.Email)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
