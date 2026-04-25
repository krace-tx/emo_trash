package chatlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/chat/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/chat/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAIHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAIHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAIHistoryLogic {
	return &GetAIHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取 AI 对话历史
func (l *GetAIHistoryLogic) GetAIHistory(in *pb.GetAIHistoryReq) (*pb.GetAIHistoryResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetAIHistoryResp{}, nil
}
