package articlelogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/article/article"
	"github.com/krace-tx/emo_trash/app/rpc/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryPendingArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryPendingArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryPendingArticlesLogic {
	return &QueryPendingArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询待审核的文章列表
func (l *QueryPendingArticlesLogic) QueryPendingArticles(in *article.QueryPendingArticlesRequest) (*article.QueryPendingArticlesResponse, error) {
	// todo: add your logic here and delete this line

	return &article.QueryPendingArticlesResponse{}, nil
}
