package authlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginByMobileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginByMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByMobileLogic {
	return &LoginByMobileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户登录
func (l *LoginByMobileLogic) LoginByMobile(in *pb.LoginByMobileReq) (*pb.LoginResp, error) {
	// todo: add your logic here and delete this line

	return &pb.LoginResp{}, nil
}
