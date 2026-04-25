package post

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"
	"github.com/krace-tx/emo_trash/app/rpc/post/client/post"
	consts "github.com/krace-tx/emo_trash/pkg/constant"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePostLogic) CreatePost(req *types.CreatePostReq) (resp *types.CommonResp, err error) {
	userId := l.ctx.Value(consts.UserId).(string)

	data, err := l.svcCtx.Post.CreatePost(l.ctx, &post.CreatePostReq{
		AuthorId:    userId,
		Title:       req.Title,
		Content:     req.Content,
		Images:      req.Images,
		IsAnonymous: req.IsAnonymous,
	})
	if err != nil {
		l.Logger.Errorf("发布帖子失败: %v, user_id=%s", err, userId)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
