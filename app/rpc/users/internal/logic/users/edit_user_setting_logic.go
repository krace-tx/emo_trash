package userslogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/users/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/users/users"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditUserSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEditUserSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditUserSettingLogic {
	return &EditUserSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 编辑隐私权限
func (l *EditUserSettingLogic) EditUserSetting(in *users.EditUserSettingRequest) (*users.EditUserSettingResponse, error) {
	// todo: add your logic here and delete this line

	return &users.EditUserSettingResponse{}, nil
}
