package post

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"
	"github.com/krace-tx/emo_trash/app/rpc/post/client/post"
	consts "github.com/krace-tx/emo_trash/pkg/constant"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePostLogic {
	return &DeletePostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePostLogic) DeletePost(req *types.DeletePostReq) (resp *types.CommonResp, err error) {
	userId := l.ctx.Value(consts.UserId).(string)

	data, err := l.svcCtx.Post.DeletePost(l.ctx, &post.DeletePostReq{
		Id:     req.Id,
		UserId: userId,
	})
	if err != nil {
		l.Logger.Errorf("閸掔娀娅庣敮鏍х摍婢惰精瑙? %v, user_id=%s, post_id=%s", err, userId, req.Id)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
