package authlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindThirdPartyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindThirdPartyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindThirdPartyLogic {
	return &BindThirdPartyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 绑定第三方登录
func (l *BindThirdPartyLogic) BindThirdParty(in *pb.BindThirdPartyReq) (*pb.BindThirdPartyResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BindThirdPartyResp{}, nil
}
