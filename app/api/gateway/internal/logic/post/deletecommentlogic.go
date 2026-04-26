// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package post

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"
	"github.com/krace-tx/emo_trash/app/rpc/post/client/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCommentLogic) DeleteComment(req *types.DeleteCommentReq) (resp *types.CommonResp, err error) {
	userId := l.ctx.Value("user_id").(string)

	_, err = l.svcCtx.Post.DeleteComment(l.ctx, &post.DeleteCommentReq{
		Id:     req.Id,
		UserId: userId,
	})
	if err != nil {
		l.Logger.Errorf("删除回声失败: %v, comment_id=%s, user_id=%s", err, req.Id, userId)
		return types.Error(err), nil
	}

	return types.Success(nil), nil
}
