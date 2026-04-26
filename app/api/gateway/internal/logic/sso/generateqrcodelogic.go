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

type GenerateQrcodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateQrcodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateQrcodeLogic {
	return &GenerateQrcodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateQrcodeLogic) GenerateQrcode(req *types.GenerateQrcodeReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.Sso.GenerateQrcode(l.ctx, &auth.GenerateQrcodeReq{
		DeviceId: req.DeviceId,
	})
	if err != nil {
		l.Logger.Errorf("生成二维码失败: %v, device_id=%s", err, req.DeviceId)
		return types.Error(err), nil
	}

	return types.Success(types.GenerateQrcodeResp{
		Qid:      data.Qid,
		ImageUrl: data.ImageUrl,
	}), nil
}
