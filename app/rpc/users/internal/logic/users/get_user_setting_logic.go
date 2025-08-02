package userslogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/users/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/users/users"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserSettingLogic {
	return &GetUserSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取隐私权限信息
func (l *GetUserSettingLogic) GetUserSetting(in *users.GetUserSettingRequest) (*users.GetUserSettingResponse, error) {
	// todo: add your logic here and delete this line

	return &users.GetUserSettingResponse{}, nil
}
