package sociallogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/users/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/users/users"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMutualFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMutualFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMutualFriendsLogic {
	return &GetMutualFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询互关朋友
func (l *GetMutualFriendsLogic) GetMutualFriends(in *users.GetMutualFriendsRequest) (*users.GetMutualFriendsResponse, error) {
	// todo: add your logic here and delete this line

	return &users.GetMutualFriendsResponse{}, nil
}
