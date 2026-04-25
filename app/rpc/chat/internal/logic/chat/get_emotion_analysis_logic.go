package chatlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/chat/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/chat/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmotionAnalysisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmotionAnalysisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmotionAnalysisLogic {
	return &GetEmotionAnalysisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取本次对话的情感分析结果（供前端动态 UI 使用）
func (l *GetEmotionAnalysisLogic) GetEmotionAnalysis(in *pb.GetEmotionAnalysisReq) (*pb.GetEmotionAnalysisResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetEmotionAnalysisResp{}, nil
}
