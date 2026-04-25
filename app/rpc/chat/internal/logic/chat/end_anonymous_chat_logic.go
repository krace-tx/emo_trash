package chatlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/chat/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/chat/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type EndAnonymousChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEndAnonymousChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EndAnonymousChatLogic {
	return &EndAnonymousChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 结束匿名聊天（双方任意一方可结束）
func (l *EndAnonymousChatLogic) EndAnonymousChat(in *pb.EndAnonymousChatReq) (*pb.CommonResp, error) {
	// todo: add your logic here and delete this line

	return &pb.CommonResp{}, nil
}
