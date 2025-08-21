package sso

import (
	"context"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindMobileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBindMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindMobileLogic {
	return &BindMobileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindMobileLogic) BindMobile(req *types.BindMobileReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.Auth.BindMobile(l.ctx, &auth.BindMobileReq{
		Mobile:  req.Mobile,
		SmsCode: req.SmsCode,
	})
	if err != nil {
		l.Logger.Errorf("BindMobile failed, err: %v", err)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
