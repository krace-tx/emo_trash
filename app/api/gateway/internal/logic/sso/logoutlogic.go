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
		Token:      req.Token,
		DeviceType: req.Device_type,
	})
	if err != nil {
		l.Logger.Errorf("Logout failed, err: %v", err)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
