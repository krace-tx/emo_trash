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

type GetComfortMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetComfortMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetComfortMessageLogic {
	return &GetComfortMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetComfortMessageLogic) GetComfortMessage(req *types.GetComfortMessageReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.Sso.GetComfortMessage(l.ctx, &auth.GetComfortMessageReq{})
	if err != nil {
		l.Logger.Errorf("获取治愈文案失败: %v", err)
		return types.Error(err), nil
	}

	return types.Success(types.GetComfortMessageResp{
		Title:           data.Title,
		Subtitle:        data.Subtitle,
		ButtonText:      data.ButtonText,
		IllustrationUrl: data.IllustrationUrl,
	}), nil
}
