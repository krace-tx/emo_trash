package sso

import (
	"context"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbindEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnbindEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindEmailLogic {
	return &UnbindEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnbindEmailLogic) UnbindEmail(req *types.UnbindEmailReq) (resp *types.CommonResp, err error) {

	data, err := l.svcCtx.Auth.UnbindEmail(l.ctx, &auth.UnbindEmailReq{
		Email:     req.Email,
		EmailCode: req.EmailCode,
	})
	if err != nil {
		l.Logger.Errorf("UnbindEmail failed, err: %v", err)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
