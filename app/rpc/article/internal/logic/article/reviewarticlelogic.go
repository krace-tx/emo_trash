package articlelogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/article/article"
	"github.com/krace-tx/emo_trash/app/rpc/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReviewArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReviewArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReviewArticleLogic {
	return &ReviewArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 审核文章（通过/未通过）
func (l *ReviewArticleLogic) ReviewArticle(in *article.ReviewArticleRequest) (*article.ReviewArticleResponse, error) {
	// todo: add your logic here and delete this line

	return &article.ReviewArticleResponse{}, nil
}
