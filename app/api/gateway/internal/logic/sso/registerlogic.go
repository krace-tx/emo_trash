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
	data, err := l.svcCtx.Auth.Register(l.ctx, &auth.RegisterReq{
		Mobile:   req.Mobile,
		Password: req.Password,
		Account:  req.Account,
		SmsCode:  req.SmsCode,
	})

	if err != nil {
		l.Logger.Errorf("Register failed, err: %v", err)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
