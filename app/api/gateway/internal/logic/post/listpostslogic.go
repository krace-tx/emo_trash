package post

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"
	"github.com/krace-tx/emo_trash/app/rpc/post/client/post"
	consts "github.com/krace-tx/emo_trash/pkg/constant"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPostsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPostsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPostsLogic {
	return &ListPostsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPostsLogic) ListPosts(req *types.ListPostsReq) (resp *types.CommonResp, err error) {
	userId, _ := l.ctx.Value(consts.UserId).(string)

	data, err := l.svcCtx.Post.ListPosts(l.ctx, &post.ListPostsReq{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: int32(req.PageSize),
		Type:     req.Type,
	})
	if err != nil {
		l.Logger.Errorf("閼惧嘲褰囩敮鏍х摍閸掓銆冩径杈Е: %v, user_id=%s", err, userId)
		return types.Error(err), nil
	}

	list := make([]types.PostInfo, 0, len(data.List))
	for _, p := range data.List {
		list = append(list, mapPostInfo(p))
	}

	return types.Success(&types.ListPostsResp{
		List:       list,
		NextCursor: data.NextCursor,
		HasMore:    data.HasMore,
	}), nil
}
