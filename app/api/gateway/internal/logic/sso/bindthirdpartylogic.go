package sso

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindThirdPartyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBindThirdPartyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindThirdPartyLogic {
	return &BindThirdPartyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindThirdPartyLogic) BindThirdParty(req *types.BindThirdPartyReq) (resp *types.CommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
