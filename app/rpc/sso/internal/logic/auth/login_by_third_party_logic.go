package authlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginByThirdPartyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginByThirdPartyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByThirdPartyLogic {
	return &LoginByThirdPartyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 三方登录
func (l *LoginByThirdPartyLogic) LoginByThirdParty(in *pb.LoginByThirdPartyReq) (*pb.LoginResp, error) {
	// TODO: 实现三方登录逻辑 (例如微信、QQ OAuth2 流程)
	// 1. 调用三方平台 API 验证 Code
	// 2. 获取 OpenID/UnionID
	// 3. 查找或创建用户 (绑定平台账号)
	// 4. 生成 JWT Token

	l.Logger.Infof("收到三方登录请求: platform=%s, code=%s", in.Platform, in.Code)

	return &pb.LoginResp{
		AccessToken: "mock_access_token_for_" + in.Platform,
	}, nil
}
