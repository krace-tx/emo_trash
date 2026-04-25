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
	data, err := l.svcCtx.Sso.ResetPassword(l.ctx, &auth.ResetPasswordReq{
		Email:       req.Email,
		EmailCode:   req.EmailCode,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		l.Logger.Errorf("重置密码失败: %v, email=%s", err, req.Email)
		return types.Error(err), nil
	}
	if data.GetSuccess() {
		return types.SuccessWithMsg(nil, data.GetMessage()), nil
	}

	return types.ParamError(data.GetMessage()), nil
}
