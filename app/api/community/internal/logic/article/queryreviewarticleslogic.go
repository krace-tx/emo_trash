package article

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryReviewArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询审核通过的文章列表
func NewQueryReviewArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryReviewArticlesLogic {
	return &QueryReviewArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryReviewArticlesLogic) QueryReviewArticles(req *types.QueryReviewArticlesReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
