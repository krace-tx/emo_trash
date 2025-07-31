package sociallogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/users/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/users/users"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowerListLogic {
	return &GetFollowerListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户的粉丝列表
func (l *GetFollowerListLogic) GetFollowerList(in *users.GetFollowerListRequest) (*users.GetFollowerListResponse, error) {
	// todo: add your logic here and delete this line

	return &users.GetFollowerListResponse{}, nil
}
