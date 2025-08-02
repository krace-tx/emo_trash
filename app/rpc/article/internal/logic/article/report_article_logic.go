package articlelogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/article/article"
	"github.com/krace-tx/emo_trash/app/rpc/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReportArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportArticleLogic {
	return &ReportArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 举报文章
func (l *ReportArticleLogic) ReportArticle(in *article.ReportArticleRequest) (*article.ReportArticleResponse, error) {
	// todo: add your logic here and delete this line

	return &article.ReportArticleResponse{}, nil
}
