package authlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbindEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnbindEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindEmailLogic {
	return &UnbindEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 解绑邮箱
func (l *UnbindEmailLogic) UnbindEmail(in *pb.UnbindEmailReq) (*pb.UnbindEmailResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UnbindEmailResp{}, nil
}
