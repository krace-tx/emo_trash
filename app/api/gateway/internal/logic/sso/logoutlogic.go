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

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.Auth.Logout(l.ctx, &auth.LogoutReq{
		Token: req.Token,
	})
	if err != nil {
		l.Logger.Errorf("登出失败: %v", err)
		return types.Error(err), nil
	}
	if data.GetSuccess() {
		return types.SuccessWithMsg(nil, data.GetMessage()), nil
	}

	return types.ParamError(data.GetMessage()), nil
}
