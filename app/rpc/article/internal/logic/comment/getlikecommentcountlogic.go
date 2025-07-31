package commentlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/article/article"
	"github.com/krace-tx/emo_trash/app/rpc/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLikeCommentCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLikeCommentCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLikeCommentCountLogic {
	return &GetLikeCommentCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取评论点赞的数量
func (l *GetLikeCommentCountLogic) GetLikeCommentCount(in *article.GetLikeCommentCountRequest) (*article.GetLikeCommentCountResponse, error) {
	// todo: add your logic here and delete this line

	return &article.GetLikeCommentCountResponse{}, nil
}
