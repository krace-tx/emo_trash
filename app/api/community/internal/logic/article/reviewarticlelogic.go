package article

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReviewArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 审核文章（通过/未通过）
func NewReviewArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReviewArticleLogic {
	return &ReviewArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReviewArticleLogic) ReviewArticle(req *types.ReviewArticleReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
