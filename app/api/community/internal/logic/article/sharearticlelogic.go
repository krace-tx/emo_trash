package article

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分享文章
func NewShareArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareArticleLogic {
	return &ShareArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareArticleLogic) ShareArticle(req *types.ShareArticleReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
