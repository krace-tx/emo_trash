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

func (l *ConfirmQrcodeLoginLogic) ConfirmQrcodeLogin(req *types.ConfirmQrcodeLoginReq) (resp *types.CommonResp, err error) {
	userId := l.ctx.Value("user_id").(string)

	_, err = l.svcCtx.Sso.ConfirmQrcodeLogin(l.ctx, &auth.ConfirmQrcodeLoginReq{
		Qid:    req.Qid,
		UserId: userId,
	})
	if err != nil {
		l.Logger.Errorf("确认扫码登录失败: %v, qid=%s", err, req.Qid)
		return types.Error(err), nil
	}

	return types.Success(nil), nil
}
