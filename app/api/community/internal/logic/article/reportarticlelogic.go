package article

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 举报文章
func NewReportArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportArticleLogic {
	return &ReportArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReportArticleLogic) ReportArticle(req *types.ReportArticleReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
