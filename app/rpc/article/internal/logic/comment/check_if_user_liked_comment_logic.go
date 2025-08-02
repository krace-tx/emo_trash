package commentlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/article/article"
	"github.com/krace-tx/emo_trash/app/rpc/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckIfUserLikedCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckIfUserLikedCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckIfUserLikedCommentLogic {
	return &CheckIfUserLikedCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckIfUserLikedCommentLogic) CheckIfUserLikedComment(in *article.CheckIfUserLikedCommentRequest) (*article.CheckIfUserLikedCommentResponse, error) {
	// todo: add your logic here and delete this line

	return &article.CheckIfUserLikedCommentResponse{}, nil
}
