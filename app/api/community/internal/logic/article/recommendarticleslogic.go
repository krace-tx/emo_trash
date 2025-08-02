package article

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecommendArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 推荐文章列表
func NewRecommendArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecommendArticlesLogic {
	return &RecommendArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecommendArticlesLogic) RecommendArticles(req *types.RecommendArticlesReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
