package sociallogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/users/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/users/users"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBlockedUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBlockedUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBlockedUsersLogic {
	return &GetBlockedUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取屏蔽用户列表
func (l *GetBlockedUsersLogic) GetBlockedUsers(in *users.GetBlockedUsersRequest) (*users.GetBlockedUsersResponse, error) {
	// todo: add your logic here and delete this line

	return &users.GetBlockedUsersResponse{}, nil
}
