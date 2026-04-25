package chatlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/chat/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/chat/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMatchStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMatchStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMatchStatusLogic {
	return &GetMatchStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询当前匹配状态
func (l *GetMatchStatusLogic) GetMatchStatus(in *pb.GetMatchStatusReq) (*pb.GetMatchStatusResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetMatchStatusResp{}, nil
}
