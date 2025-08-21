package sso

import (
	"context"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"

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

	data, err := l.svcCtx.Auth.UnbindThirdParty(l.ctx, &auth.UnbindThirdPartyReq{
		Platform: req.Platform,
		OpenId:   req.OpenId,
		UnionId:  req.UnionId,
	})
	if err != nil {
		l.Logger.Errorf("UnbindThirdParty failed, err: %v", err)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
