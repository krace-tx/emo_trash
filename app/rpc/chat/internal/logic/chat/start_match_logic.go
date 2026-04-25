package chatlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/chat/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/chat/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type StartMatchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStartMatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartMatchLogic {
	return &StartMatchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ---------- 匿名匹配 ----------
func (l *StartMatchLogic) StartMatch(in *pb.StartMatchReq) (*pb.StartMatchResp, error) {
	// todo: add your logic here and delete this line

	return &pb.StartMatchResp{}, nil
}
