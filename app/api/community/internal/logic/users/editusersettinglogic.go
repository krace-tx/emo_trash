package users

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditUserSettingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 编辑用户隐私权限
func NewEditUserSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditUserSettingLogic {
	return &EditUserSettingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditUserSettingLogic) EditUserSetting(req *types.EditUserSettingReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
