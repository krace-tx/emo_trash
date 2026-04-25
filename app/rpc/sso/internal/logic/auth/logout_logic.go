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

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 登出
func (l *LogoutLogic) Logout(in *pb.LogoutReq) (*pb.CommonResp, error) {
	// 1. 解析 Token 获取过期时间
	claims, err := authx.ParseJwtToken(in.Token, l.svcCtx.Config.JWT.AccessSecret)
	if err != nil {
		// 已经失效或错误的 Token 直接返回成功
		return &pb.CommonResp{Success: true, Message: "Logout success"}, nil
	}

	// 2. 计算剩余有效期并加入黑名单
	exp, ok := claims[consts.Expire].(float64)
	if !ok {
		return nil, errx.ErrAuthForbidden
	}

	remaining := time.Unix(int64(exp), 0).Sub(time.Now())
	if remaining > 0 {
		if err := l.addToBlacklist(in.Token, remaining); err != nil {
			l.Logger.Errorf("加入黑名单失败: %v, token=%s", err, in.Token[:10])
			return nil, errx.ErrAuthTokenBlacklistFail
		}
	}

	l.Logger.Infof("用户登出成功: user_id=%v", claims[consts.UserId])
	return &pb.CommonResp{Success: true, Message: "Logout success"}, nil
}

func (l *LogoutLogic) addToBlacklist(token string, expiration time.Duration) error {
	blacklistKey := redis.GenerateKey("sso", "token", "blacklist", token)
	return l.svcCtx.Redis.Set(blacklistKey, "1", expiration)
}
