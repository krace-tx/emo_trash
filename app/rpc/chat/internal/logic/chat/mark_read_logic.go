package chatlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/chat/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/chat/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkReadLogic {
	return &MarkReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 标记会话已读
func (l *MarkReadLogic) MarkRead(in *pb.MarkReadReq) (*pb.CommonResp, error) {
	// todo: add your logic here and delete this line

	return &pb.CommonResp{}, nil
}
