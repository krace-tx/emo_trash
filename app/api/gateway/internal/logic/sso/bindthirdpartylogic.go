package sso

import (
	"context"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"

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
	data, err := l.svcCtx.Auth.BindThirdParty(l.ctx, &auth.BindThirdPartyReq{
		Platform: req.Platform,
		OpenId:   req.OpenId,
		UnionId:  req.UnionId,
	})
	if err != nil {
		l.Logger.Errorf("BindThirdParty failed, err: %v", err)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
