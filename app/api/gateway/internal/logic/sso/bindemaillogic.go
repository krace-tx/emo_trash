package sso

import (
	"context"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBindEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindEmailLogic {
	return &BindEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindEmailLogic) BindEmail(req *types.BindEmailReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.Auth.BindEmail(l.ctx, &auth.BindEmailReq{
		Email:     req.Email,
		EmailCode: req.EmailCode,
	})
	if err != nil {
		l.Logger.Errorf("BindEmail failed, err: %v", err)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
