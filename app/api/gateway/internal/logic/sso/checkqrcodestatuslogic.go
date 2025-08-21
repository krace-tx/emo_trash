package sso

import (
	"context"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckQrcodeStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckQrcodeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckQrcodeStatusLogic {
	return &CheckQrcodeStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckQrcodeStatusLogic) CheckQrcodeStatus(req *types.QrcodeStatusReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.Auth.CheckQrcodeStatus(l.ctx, &auth.QrcodeStatusReq{
		Qid: req.Qid,
	})
	if err != nil {
		l.Logger.Errorf("CheckQrcodeStatus failed, err: %v", err)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
