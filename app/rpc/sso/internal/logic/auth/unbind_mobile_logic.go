package authlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbindMobileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnbindMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindMobileLogic {
	return &UnbindMobileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 解绑手机号
func (l *UnbindMobileLogic) UnbindMobile(in *pb.UnbindMobileReq) (*pb.UnbindMobileResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UnbindMobileResp{}, nil
}
