package article

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActionArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 编辑文章活动状态
func NewActionArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActionArticleLogic {
	return &ActionArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ActionArticleLogic) ActionArticle(req *types.ActionArticleReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
