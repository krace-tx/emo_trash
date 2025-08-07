package authlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckQrcodeStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckQrcodeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckQrcodeStatusLogic {
	return &CheckQrcodeStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 检查二维码状态(PC端)
func (l *CheckQrcodeStatusLogic) CheckQrcodeStatus(in *pb.QrcodeStatusReq) (*pb.QrcodeStatusResp, error) {
	// todo: add your logic here and delete this line

	return &pb.QrcodeStatusResp{}, nil
}
