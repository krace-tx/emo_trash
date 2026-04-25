package chatlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/chat/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/chat/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecallMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecallMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecallMessageLogic {
	return &RecallMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 撤回消息（仅限 2 分钟内）
func (l *RecallMessageLogic) RecallMessage(in *pb.RecallMessageReq) (*pb.CommonResp, error) {
	// todo: add your logic here and delete this line

	return &pb.CommonResp{}, nil
}
