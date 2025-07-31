package articlelogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/article/article"
	"github.com/krace-tx/emo_trash/app/rpc/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveDraftArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveDraftArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveDraftArticleLogic {
	return &SaveDraftArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存草稿
func (l *SaveDraftArticleLogic) SaveDraftArticle(in *article.SaveDraftArticleRequest) (*article.SaveDraftArticleResponse, error) {
	// todo: add your logic here and delete this line

	return &article.SaveDraftArticleResponse{}, nil
}
