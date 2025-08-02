package comment

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReplyCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取回复评论
func NewGetReplyCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReplyCommentListLogic {
	return &GetReplyCommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetReplyCommentListLogic) GetReplyCommentList(req *types.GetReplyCommentListReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
