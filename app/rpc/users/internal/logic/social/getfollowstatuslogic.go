package sociallogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/users/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/users/users"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowStatusLogic {
	return &GetFollowStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询其他用户与自己的关系状态码
func (l *GetFollowStatusLogic) GetFollowStatus(in *users.GetFollowStatusRequest) (*users.GetFollowStatusResponse, error) {
	// todo: add your logic here and delete this line

	return &users.GetFollowStatusResponse{}, nil
}
