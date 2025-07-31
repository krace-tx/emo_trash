package article

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveDraftArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 保存草稿
func NewSaveDraftArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveDraftArticleLogic {
	return &SaveDraftArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveDraftArticleLogic) SaveDraftArticle(req *types.CreateArticleReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
