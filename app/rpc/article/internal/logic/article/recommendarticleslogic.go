package articlelogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/article/article"
	"github.com/krace-tx/emo_trash/app/rpc/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecommendArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecommendArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecommendArticlesLogic {
	return &RecommendArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 推荐文章列表
func (l *RecommendArticlesLogic) RecommendArticles(in *article.RecommendArticlesRequest) (*article.RecommendArticlesResponse, error) {
	// todo: add your logic here and delete this line

	return &article.RecommendArticlesResponse{}, nil
}
