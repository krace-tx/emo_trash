package sso

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbindThirdPartyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnbindThirdPartyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindThirdPartyLogic {
	return &UnbindThirdPartyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnbindThirdPartyLogic) UnbindThirdParty(req *types.UnbindThirdPartyReq) (resp *types.CommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
