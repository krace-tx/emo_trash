package sso

import (
	"context"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbindMobileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnbindMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindMobileLogic {
	return &UnbindMobileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnbindMobileLogic) UnbindMobile(req *types.UnbindMobileReq) (resp *types.CommonResp, err error) {

	data, err := l.svcCtx.Auth.UnbindMobile(l.ctx, &auth.UnbindMobileReq{
		Mobile:  req.Mobile,
		SmsCode: req.SmsCode,
	})
	if err != nil {
		l.Logger.Errorf("UnbindMobile failed, err: %v", err)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
