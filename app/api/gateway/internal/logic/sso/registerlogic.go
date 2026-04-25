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

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.Sso.Register(l.ctx, &auth.RegisterReq{
		Email:     req.Email,
		EmailCode: req.EmailCode,
		Password:  req.Password,
	})
	if err != nil {
		l.Logger.Errorf("邮箱注册失败: %v, email=%s", err, req.Email)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
