package article

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryPendingArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询待审核的文章列表
func NewQueryPendingArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryPendingArticlesLogic {
	return &QueryPendingArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryPendingArticlesLogic) QueryPendingArticles(req *types.QueryPendingArticlesReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
