package articlelogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/article/article"
	"github.com/krace-tx/emo_trash/app/rpc/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckIfUserActionArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckIfUserActionArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckIfUserActionArticleLogic {
	return &CheckIfUserActionArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询用户是否对文章进行点赞
func (l *CheckIfUserActionArticleLogic) CheckIfUserActionArticle(in *article.CheckIfUserActionArticleRequest) (*article.CheckIfUserActionArticleResponse, error) {
	// todo: add your logic here and delete this line

	return &article.CheckIfUserActionArticleResponse{}, nil
}
