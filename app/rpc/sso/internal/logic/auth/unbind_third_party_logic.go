package authlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbindThirdPartyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnbindThirdPartyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindThirdPartyLogic {
	return &UnbindThirdPartyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 解绑第三方登录
func (l *UnbindThirdPartyLogic) UnbindThirdParty(in *pb.UnbindThirdPartyReq) (*pb.UnbindThirdPartyResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UnbindThirdPartyResp{}, nil
}
