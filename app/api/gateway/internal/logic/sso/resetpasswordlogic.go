package sso

import (
	"context"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordReq) (resp *types.CommonResp, err error) {

	data, err := l.svcCtx.Auth.ResetPassword(l.ctx, &auth.ResetPasswordReq{
		Mobile:      req.Mobile,
		SmsCode:     req.SmsCode,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		l.Logger.Errorf("ResetPassword failed, err: %v", err)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
