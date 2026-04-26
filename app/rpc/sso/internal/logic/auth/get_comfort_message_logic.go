package authlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetComfortMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetComfortMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetComfortMessageLogic {
	return &GetComfortMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取每日温柔文案
func (l *GetComfortMessageLogic) GetComfortMessage(in *pb.GetComfortMessageReq) (*pb.GetComfortMessageResp, error) {
	// 在实际项目中，这里可以对接运营后台的配置，或者根据日期随机下发
	// 目前先返回默认文案对齐前端 UI
	return &pb.GetComfortMessageResp{
		Title:           "你不是一个人",
		Subtitle:        "今天也有人和你一样",
		ButtonText:      "给自己一点温柔",
		IllustrationUrl: "/static/mine/illustration.png",
	}, nil
}
