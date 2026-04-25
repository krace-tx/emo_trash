package authlogic

import (
	"context"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"
	authx "github.com/krace-tx/emo_trash/pkg/auth"
	consts "github.com/krace-tx/emo_trash/pkg/constant"
	"github.com/krace-tx/emo_trash/pkg/datastore/redis"
	errx "github.com/krace-tx/emo_trash/pkg/err"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyTokenLogic {
	return &VerifyTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 验证会话
func (l *VerifyTokenLogic) VerifyToken(in *pb.VerifyReq) (*pb.VerifyResp, error) {
	// 1. 解析并验证 JWT 签名与过期
	claims, err := authx.ParseJwtToken(in.Token, l.svcCtx.Config.JWT.AccessSecret)
	if err != nil {
		l.Logger.Errorf("JWT 解析失败: %v", err)
		return nil, errx.ErrAuthForbidden
	}

	// 2. 检查黑名单
	if l.isBlacklisted(in.Token) {
		l.Logger.Infof("尝试使用已拉黑的令牌: token=%s", in.Token[:10]+"...")
		return nil, errx.ErrAuthForbidden
	}

	userId, ok := claims[consts.UserId].(string)
	if !ok {
		return nil, errx.ErrAuthForbidden
	}
	email, _ := claims["email"].(string)

	return &pb.VerifyResp{
		UserId: userId,
		Email:  email,
	}, nil
}

func (l *VerifyTokenLogic) isBlacklisted(token string) bool {
	blacklistKey := redis.GenerateKey("sso", "token", "blacklist", token)
	var blacklisted string
	_ = l.svcCtx.Redis.Get(blacklistKey, &blacklisted)
	return blacklisted != ""
}
