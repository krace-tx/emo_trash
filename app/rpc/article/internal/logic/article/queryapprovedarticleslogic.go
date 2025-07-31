package articlelogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/article/article"
	"github.com/krace-tx/emo_trash/app/rpc/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryApprovedArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryApprovedArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryApprovedArticlesLogic {
	return &QueryApprovedArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询审核通过的文章列表
func (l *QueryApprovedArticlesLogic) QueryApprovedArticles(in *article.QueryApprovedArticlesRequest) (*article.QueryApprovedArticlesResponse, error) {
	// todo: add your logic here and delete this line

	return &article.QueryApprovedArticlesResponse{}, nil
}
