package post

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/gateway/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/gateway/internal/types"
	"github.com/krace-tx/emo_trash/app/rpc/post/client/post"
	consts "github.com/krace-tx/emo_trash/pkg/constant"

	"github.com/zeromicro/go-zero/core/logx"
)

type StarPostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStarPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StarPostLogic {
	return &StarPostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StarPostLogic) StarPost(req *types.StarPostReq) (resp *types.CommonResp, err error) {
	userId := l.ctx.Value(consts.UserId).(string)

	data, err := l.svcCtx.Post.StarPost(l.ctx, &post.StarPostReq{
		Id:     req.Id,
		UserId: userId,
	})
	if err != nil {
		l.Logger.Errorf("閺€鎯版鐢牕鐡欐径杈Е: %v, user_id=%s, post_id=%s", err, userId, req.Id)
		return types.Error(err), nil
	}

	return types.Success(data), nil
}
