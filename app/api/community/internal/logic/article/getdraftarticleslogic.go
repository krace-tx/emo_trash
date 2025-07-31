package article

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDraftArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查看用户的草稿列表
func NewGetDraftArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDraftArticlesLogic {
	return &GetDraftArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDraftArticlesLogic) GetDraftArticles(req *types.GetDraftArticlesReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
