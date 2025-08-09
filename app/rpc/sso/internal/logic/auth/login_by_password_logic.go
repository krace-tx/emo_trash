package authlogic

import (
	"context"
	"errors"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginByPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginByPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByPasswordLogic {
	return &LoginByPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 账号密码登录
func (l *LoginByPasswordLogic) LoginByPassword(in *pb.LoginByPasswordReq) (*pb.LoginResp, error) {
	if in.Account == "" || in.Password == "" {
		return nil, errors.New("账号或密码不能为空")
	}

	//engine := rdb.NewEngine[*model.UserAuth](rdb.M)

	return &pb.LoginResp{}, nil
}
