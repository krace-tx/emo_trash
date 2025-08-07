package authlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateQrcodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateQrcodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateQrcodeLogic {
	return &GenerateQrcodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 生成登录二维码(PC端)
func (l *GenerateQrcodeLogic) GenerateQrcode(in *pb.QrcodeReq) (*pb.QrcodeResp, error) {
	// todo: add your logic here and delete this line

	return &pb.QrcodeResp{}, nil
}
