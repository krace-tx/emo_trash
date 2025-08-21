package sso

import (
	"context"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordByEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPasswordByEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordByEmailLogic {
	return &ResetPasswordByEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordByEmailLogic) ResetPasswordByEmail(req *types.ResetPasswordByEmailReq) (resp *types.CommonResp, err error) {

	data, err := l.svcCtx.Auth.ResetPasswordByEmail(l.ctx, &auth.ResetPasswordByEmailReq{
		Email:       req.Email,
		EmailCode:   req.EmailCode,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		l.Logger.Errorf("ResetPasswordByEmail failed, err: %v", err)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
