package comment

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 回复评论信息
func NewReplyCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyCommentLogic {
	return &ReplyCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReplyCommentLogic) ReplyComment(req *types.ReplyCommentReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
