package commentlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/article/article"
	"github.com/krace-tx/emo_trash/app/rpc/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReplyCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetReplyCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReplyCommentLogic {
	return &GetReplyCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取回复评论
func (l *GetReplyCommentLogic) GetReplyComment(in *article.GetReplyCommentRequest) (*article.GetReplyCommentResponse, error) {
	// todo: add your logic here and delete this line

	return &article.GetReplyCommentResponse{}, nil
}
