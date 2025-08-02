package users

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 编辑用户信息
func NewEditUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditUserInfoLogic {
	return &EditUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditUserInfoLogic) EditUserInfo(req *types.EditUserInfoReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
