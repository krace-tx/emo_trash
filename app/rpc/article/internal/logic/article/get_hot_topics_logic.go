package articlelogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/article/article"
	"github.com/krace-tx/emo_trash/app/rpc/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHotTopicsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetHotTopicsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHotTopicsLogic {
	return &GetHotTopicsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取热门帖子
func (l *GetHotTopicsLogic) GetHotTopics(in *article.GetHotTopicsRequest) (*article.GetHotTopicsResponse, error) {
	// todo: add your logic here and delete this line

	return &article.GetHotTopicsResponse{}, nil
}
