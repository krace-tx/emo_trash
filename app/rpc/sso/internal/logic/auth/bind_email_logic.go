package authlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindEmailLogic {
	return &BindEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 绑定邮箱
func (l *BindEmailLogic) BindEmail(in *pb.BindEmailReq) (*pb.BindEmailResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BindEmailResp{}, nil
}
