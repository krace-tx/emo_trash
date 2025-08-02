package authlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/sso"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmQrcodeLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfirmQrcodeLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmQrcodeLoginLogic {
	return &ConfirmQrcodeLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 手机端确认登录
func (l *ConfirmQrcodeLoginLogic) ConfirmQrcodeLogin(in *sso.QrcodeConfirmReq) (*sso.LoginResp, error) {
	// todo: add your logic here and delete this line

	return &sso.LoginResp{}, nil
}
