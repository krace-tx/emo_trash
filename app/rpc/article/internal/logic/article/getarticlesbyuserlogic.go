package articlelogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/article/article"
	"github.com/krace-tx/emo_trash/app/rpc/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticlesByUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArticlesByUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticlesByUserLogic {
	return &GetArticlesByUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查看用户发布的文章列表
func (l *GetArticlesByUserLogic) GetArticlesByUser(in *article.GetArticlesByUserRequest) (*article.GetArticlesByUserResponse, error) {
	// todo: add your logic here and delete this line

	return &article.GetArticlesByUserResponse{}, nil
}
