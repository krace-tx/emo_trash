package article

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户收藏的文章
func NewGetFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteLogic {
	return &GetFavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFavoriteLogic) GetFavorite(req *types.GetFavoriteReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
