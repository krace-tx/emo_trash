package authlogic

import (
	"context"
	"time"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"
	authx "github.com/krace-tx/emo_trash/pkg/auth"
	consts "github.com/krace-tx/emo_trash/pkg/constant"
	"github.com/krace-tx/emo_trash/pkg/datastore/redis"
	errx "github.com/krace-tx/emo_trash/pkg/err"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 刷新 Token
func (l *RefreshTokenLogic) RefreshToken(in *pb.RefreshTokenReq) (*pb.LoginResp, error) {
	// 1. 解析 Refresh Token
	claims, err := authx.ParseJwtToken(in.RefreshToken, l.svcCtx.Config.JWT.RefreshSecret)
	if err != nil {
		l.Logger.Errorf("Refresh Token 解析失败: %v", err)
		return nil, errx.ErrAuthForbidden
	}

	// 2. 检查 Refresh Token 是否被拉黑
	if l.isBlacklisted(in.RefreshToken) {
		l.Logger.Infof("尝试使用已拉黑的刷新令牌: token=%s", in.RefreshToken[:10])
		return nil, errx.ErrAuthForbidden
	}

	// 3. 生成新 Token 对
	userId := claims[consts.UserId].(string)
	email := claims["email"].(string)

	newClaims := map[string]any{
		consts.UserId: userId,
		"email":       email,
	}

	accessToken, err := authx.GenJwtToken(
		l.svcCtx.Config.JWT.AccessSecret,
		l.svcCtx.Config.JWT.AccessExpire,
		newClaims,
	)
	if err != nil {
		return nil, errx.ErrAuthGenAccessTokenFail
	}

	refreshToken, err := authx.GenJwtToken(
		l.svcCtx.Config.JWT.RefreshSecret,
		l.svcCtx.Config.JWT.RefreshExpire,
		newClaims,
	)
	if err != nil {
		return nil, errx.ErrAuthGenRefreshTokenFail
	}

	// 4. 将旧 Refresh Token 拉黑（旋转机制）
	exp, _ := claims[consts.Expire].(float64)
	remaining := time.Unix(int64(exp), 0).Sub(time.Now())
	if remaining > 0 {
		_ = l.addToBlacklist(in.RefreshToken, remaining)
	}

	l.Logger.Infof("用户刷新令牌成功: user_id=%s", userId)
	return &pb.LoginResp{
		AccessToken:        accessToken,
		AccessTokenExpire:  l.svcCtx.Config.JWT.AccessExpire,
		RefreshToken:       refreshToken,
		RefreshTokenExpire: l.svcCtx.Config.JWT.RefreshExpire,
	}, nil
}

func (l *RefreshTokenLogic) isBlacklisted(token string) bool {
	blacklistKey := redis.GenerateKey("sso", "token", "blacklist", token)
	var blacklisted string
	_ = l.svcCtx.Redis.Get(blacklistKey, &blacklisted)
	return blacklisted != ""
}

func (l *RefreshTokenLogic) addToBlacklist(token string, expiration time.Duration) error {
	blacklistKey := redis.GenerateKey("sso", "token", "blacklist", token)
	return l.svcCtx.Redis.Set(blacklistKey, "1", expiration)
}
