package sso

import (
	"context"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyTokenLogic {
	return &VerifyTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyTokenLogic) VerifyToken(req *types.VerifyReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.Auth.VerifyToken(l.ctx, &auth.VerifyReq{
		Token:      req.Token,
		DeviceType: req.DeviceType,
	})
	if err != nil {
		l.Logger.Errorf("VerifyToken failed, err: %v", err)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
