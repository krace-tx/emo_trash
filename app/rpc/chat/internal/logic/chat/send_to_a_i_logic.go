package chatlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/chat/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/chat/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendToAILogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendToAILogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendToAILogic {
	return &SendToAILogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ---------- AI 情感陪伴 ----------
func (l *SendToAILogic) SendToAI(in *pb.SendToAIReq) (*pb.SendToAIResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SendToAIResp{}, nil
}
