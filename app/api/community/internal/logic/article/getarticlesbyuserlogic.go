package article

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticlesByUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查看用户发布的文章列表
func NewGetArticlesByUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticlesByUserLogic {
	return &GetArticlesByUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticlesByUserLogic) GetArticlesByUser(req *types.GetArticlesByUserReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
