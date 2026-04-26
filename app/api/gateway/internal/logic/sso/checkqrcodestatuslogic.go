// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package sso

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"
	"github.com/krace-tx/emo_trash/app/rpc/sso/client/auth"

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

func (l *CheckQrcodeStatusLogic) CheckQrcodeStatus(req *types.CheckQrcodeStatusReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.Sso.CheckQrcodeStatus(l.ctx, &auth.CheckQrcodeStatusReq{
		Qid: req.Qid,
	})
	if err != nil {
		l.Logger.Errorf("获取二维码状态失败: %v, qid=%s", err, req.Qid)
		return types.Error(err), nil
	}

	return types.Success(types.CheckQrcodeStatusResp{
		Status:       data.Status,
		Token:        data.AccessToken,
		RefreshToken: data.RefreshToken,
	}), nil
}
