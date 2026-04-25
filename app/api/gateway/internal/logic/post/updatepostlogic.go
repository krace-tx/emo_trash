package post

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"
	"github.com/krace-tx/emo_trash/app/rpc/post/client/post"
	consts "github.com/krace-tx/emo_trash/pkg/constant"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePostLogic {
	return &UpdatePostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePostLogic) UpdatePost(req *types.UpdatePostReq) (resp *types.CommonResp, err error) {
	userId := l.ctx.Value(consts.UserId).(string)

	data, err := l.svcCtx.Post.UpdatePost(l.ctx, &post.UpdatePostReq{
		Id:      req.Id,
		UserId:  userId,
		Title:   req.Title,
		Content: req.Content,
		Images:  req.Images,
	})
	if err != nil {
		l.Logger.Errorf("閺囧瓨鏌婄敮鏍х摍婢惰精瑙? %v, user_id=%s, post_id=%s", err, userId, req.Id)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
