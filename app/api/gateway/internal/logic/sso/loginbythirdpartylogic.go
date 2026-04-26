// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package sso

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginByThirdPartyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginByThirdPartyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByThirdPartyLogic {
	return &LoginByThirdPartyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginByThirdPartyLogic) LoginByThirdParty(req *types.LoginByThirdPartyReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.Sso.LoginByThirdParty(l.ctx, &auth.LoginByThirdPartyReq{
		Platform: req.Platform,
		Code:     req.Code,
	})
	if err != nil {
		l.Logger.Errorf("三方登录失败: %v, platform=%s", err, req.Platform)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
