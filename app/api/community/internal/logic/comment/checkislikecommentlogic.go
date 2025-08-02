package comment

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckIsLikeCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 检查用户是否对评论进行点赞
func NewCheckIsLikeCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckIsLikeCommentLogic {
	return &CheckIsLikeCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckIsLikeCommentLogic) CheckIsLikeComment(req *types.CheckIsLikeCommentReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
