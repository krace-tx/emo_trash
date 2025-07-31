package social

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMutualFriendsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询互关朋友
func NewGetMutualFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMutualFriendsLogic {
	return &GetMutualFriendsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMutualFriendsLogic) GetMutualFriends(req *types.GetMutualFriendsReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
