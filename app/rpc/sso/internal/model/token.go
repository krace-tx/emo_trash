package model

import (
	"context"
	"time"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	authx "github.com/krace-tx/emo_trash/pkg/auth"
	"github.com/zeromicro/go-zero/core/logx"
)

type Token struct {
	AccessToken   string `json:"access_token"`
	AccessExpire  int64  `json:"access_expire"`
	RefreshToken  string `json:"refresh_token"`
	RefreshExpire int64  `json:"refresh_expire"`
}

type TokenModel struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTokenModel(ctx context.Context, svcCtx *svc.ServiceContext) *TokenModel {
	return &TokenModel{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Generate tokens
func (t *TokenModel) GenerateTokens(auth UserAuth) (token *Token, err error) {
	token = new(Token)
	// 当前时间戳 + 过期时间
	token.AccessExpire = time.Now().Unix() + t.svcCtx.Config.JWT.AccessExpire // 秒级时间戳
	// 访问令牌过期时间
	token.AccessToken, err = authx.GenJwtToken(
		t.svcCtx.Config.JWT.AccessSecret,
		t.svcCtx.Config.JWT.AccessExpire,
		map[string]any{
			"user_id":       auth.UserID,
			"account":       auth.Account,
			"platform":      auth.Platform,
			"last_login_ip": auth.LastLoginIP,
		},
	)
	if err != nil {
		t.Logger.Error(err)
		return
	}
	// 刷新令牌过期时间
	token.RefreshExpire = time.Now().Unix() + t.svcCtx.Config.JWT.RefreshExpire // 秒级时间戳
	token.RefreshToken, err = authx.GenJwtToken(
		t.svcCtx.Config.JWT.RefreshSecret,
		t.svcCtx.Config.JWT.RefreshExpire,
		map[string]any{
			"user_id":       auth.UserID,
			"account":       auth.Account,
			"platform":      auth.Platform,
			"last_login_ip": auth.LastLoginIP,
		},
	)
	if err != nil {
		t.Logger.Error(err)
		return
	}
	return
}
