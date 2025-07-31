package sociallogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/users/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/users/users"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowListLogic {
	return &GetFollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户关注列表
func (l *GetFollowListLogic) GetFollowList(in *users.GetFollowListRequest) (*users.GetFollowListResponse, error) {
	// todo: add your logic here and delete this line

	return &users.GetFollowListResponse{}, nil
}
