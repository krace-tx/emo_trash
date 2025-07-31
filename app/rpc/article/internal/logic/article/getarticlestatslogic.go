package articlelogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/article/article"
	"github.com/krace-tx/emo_trash/app/rpc/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArticleStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleStatsLogic {
	return &GetArticleStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取文章统计信息
func (l *GetArticleStatsLogic) GetArticleStats(in *article.GetArticleStatsRequest) (*article.GetArticleStatsResponse, error) {
	// todo: add your logic here and delete this line

	return &article.GetArticleStatsResponse{}, nil
}
