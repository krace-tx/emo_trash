package article

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GlobalSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 全局搜索文章接口
func NewGlobalSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GlobalSearchLogic {
	return &GlobalSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GlobalSearchLogic) GlobalSearch(req *types.GlobalSearchReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
