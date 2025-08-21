package sso

import (
	"context"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmQrcodeLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfirmQrcodeLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmQrcodeLoginLogic {
	return &ConfirmQrcodeLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfirmQrcodeLoginLogic) ConfirmQrcodeLogin(req *types.QrcodeConfirmReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.Auth.ConfirmQrcodeLogin(l.ctx, &auth.QrcodeConfirmReq{
		Qid:      req.Qid,
		AppToken: req.AppToken,
	})
	if err != nil {
		l.Logger.Errorf("ConfirmQrcodeLogin failed, err: %v", err)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
