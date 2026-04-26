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

type ListCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCommentsLogic {
	return &ListCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListCommentsLogic) ListComments(req *types.ListCommentsReq) (resp *types.CommonResp, err error) {
	data, err := l.svcCtx.Post.ListComments(l.ctx, &post.ListCommentsReq{
		PostId:   req.PostId,
		Cursor:   req.Cursor,
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		l.Logger.Errorf("获取回声列表失败: %v, post_id=%s", err, req.PostId)
		return types.Error(err), nil
	}

	list := make([]types.CommentInfo, 0, len(data.List))
	for _, item := range data.List {
		list = append(list, types.CommentInfo{
			Id:           item.Id,
			PostId:       item.PostId,
			AuthorId:     item.AuthorId,
			AuthorName:   item.AuthorName,
			AuthorAvatar: item.AuthorAvatar,
			Content:      item.Content,
			CreatedAt:    item.CreatedAt,
		})
	}

	return types.Success(types.ListCommentsResp{
		List:       list,
		NextCursor: data.NextCursor,
		HasMore:    data.HasMore,
	}), nil
}
