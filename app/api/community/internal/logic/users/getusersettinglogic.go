package users

import (
	"context"

	"github.com/krace-tx/emo_trash/app/api/community/internal/svc"
	"github.com/krace-tx/emo_trash/app/api/community/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserSettingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户隐私权限
func NewGetUserSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserSettingLogic {
	return &GetUserSettingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserSettingLogic) GetUserSetting(req *types.GetUserSettingReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
