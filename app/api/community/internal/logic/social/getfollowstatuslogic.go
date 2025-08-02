package social

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询其他用户与自己的关系状态码
func NewGetFollowStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowStatusLogic {
	return &GetFollowStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFollowStatusLogic) GetFollowStatus(req *types.GetFollowStatusReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
