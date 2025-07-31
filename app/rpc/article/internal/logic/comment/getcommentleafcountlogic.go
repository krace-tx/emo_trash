package commentlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/article/article"
	"github.com/krace-tx/emo_trash/app/rpc/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentLeafCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentLeafCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentLeafCountLogic {
	return &GetCommentLeafCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取叶子节点的数量
func (l *GetCommentLeafCountLogic) GetCommentLeafCount(in *article.GetCommentLeafCountRequest) (*article.GetCommentLeafCountResponse, error) {
	// todo: add your logic here and delete this line

	return &article.GetCommentLeafCountResponse{}, nil
}
