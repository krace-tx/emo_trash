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

type ListMyPostsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyPostsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyPostsLogic {
	return &ListMyPostsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMyPostsLogic) ListMyPosts(req *types.ListMyPostsReq) (resp *types.CommonResp, err error) {
	userId := l.ctx.Value("user_id").(string)

	data, err := l.svcCtx.Post.ListMyPosts(l.ctx, &post.ListMyPostsReq{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		l.Logger.Errorf("获取我的情绪列表失败: %v, user_id=%s", err, userId)
		return types.Error(err), nil
	}

	list := make([]types.PostInfo, 0, len(data.List))
	for _, item := range data.List {
		list = append(list, types.PostInfo{
			Id:           item.Id,
			AuthorId:     item.AuthorId,
			AuthorName:   item.AuthorName,
			AuthorAvatar: item.AuthorAvatar,
			Title:        item.Title,
			Content:      item.Content,
			Images:       item.Images,
			Video:        item.Video,
			AiEvaluation: item.AiEvaluation,
			LikeCount:    item.LikeCount,
			CommentCount: item.CommentCount,
			StarCount:    item.StarCount,
			IsLiked:      item.IsLiked,
			IsStarred:    item.IsStarred,
			CreatedAt:    item.CreatedAt,
		})
	}

	return types.Success(types.ListPostsResp{
		List:       list,
		NextCursor: data.NextCursor,
		HasMore:    data.HasMore,
	}), nil
}
