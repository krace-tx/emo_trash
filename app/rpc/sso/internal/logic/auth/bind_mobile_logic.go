package authlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindMobileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindMobileLogic {
	return &BindMobileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 绑定手机号
func (l *BindMobileLogic) BindMobile(in *pb.BindMobileReq) (*pb.BindMobileResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BindMobileResp{}, nil
}
