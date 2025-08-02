package articlelogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/article/article"
	"github.com/krace-tx/emo_trash/app/rpc/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryRejectedArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryRejectedArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRejectedArticlesLogic {
	return &QueryRejectedArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询审核未通过的文章列表
func (l *QueryRejectedArticlesLogic) QueryRejectedArticles(in *article.QueryRejectedArticlesRequest) (*article.QueryRejectedArticlesResponse, error) {
	// todo: add your logic here and delete this line

	return &article.QueryRejectedArticlesResponse{}, nil
}
