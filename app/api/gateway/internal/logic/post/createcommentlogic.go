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

type CreateCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCommentLogic) CreateComment(req *types.CreateCommentReq) (resp *types.CommonResp, err error) {
	userId := l.ctx.Value("user_id").(string)

	_, err = l.svcCtx.Post.CreateComment(l.ctx, &post.CreateCommentReq{
		PostId:  req.PostId,
		UserId:  userId,
		Content: req.Content,
	})
	if err != nil {
		l.Logger.Errorf("发表回声失败: %v, post_id=%s, user_id=%s", err, req.PostId, userId)
		return types.Error(err), nil
	}

	return types.Success(nil), nil
}
