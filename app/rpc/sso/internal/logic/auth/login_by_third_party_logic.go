package authlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginByThirdPartyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginByThirdPartyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByThirdPartyLogic {
	return &LoginByThirdPartyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 第三方平台登录
func (l *LoginByThirdPartyLogic) LoginByThirdParty(in *pb.LoginByThirdPartyReq) (*pb.LoginResp, error) {
	// todo: add your logic here and delete this line

	return &pb.LoginResp{}, nil
}
