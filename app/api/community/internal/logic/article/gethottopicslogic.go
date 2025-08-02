package article

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHotTopicsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取热门帖子
func NewGetHotTopicsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHotTopicsLogic {
	return &GetHotTopicsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHotTopicsLogic) GetHotTopics(req *types.GetHotTopicsReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
