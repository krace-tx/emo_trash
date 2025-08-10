package authlogic

import (
	"context"
	"errors"
	"fmt" // 新增导入

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"
	model "github.com/krace-tx/emo_trash/model/user_center"
	authx "github.com/krace-tx/emo_trash/pkg/auth"
	"github.com/krace-tx/emo_trash/pkg/db/rdb"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
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
	// 1. 基础参数校验
	if in.Account == "" || in.Password == "" {
		return nil, errors.New("账号或密码不能为空")
	}

	// 2. 查询用户认证信息
	engine := rdb.NewEngine[model.UserAuth](rdb.M)
	auth, err := engine.GetByCondition(l.ctx, rdb.WithConditions("account = ?", in.Account))
	if err != nil {
		l.Logger.Errorf("查询用户认证信息失败: %v, account=%s", err, in.Account)
		return nil, errors.New("登录失败，请重试") // 避免暴露数据库错误详情
	}

	if auth == nil {
		return nil, errors.New("账号不存在")
	}

	// 3. 检查用户状态
	if auth.Status == model.UserStatusDisabled {
		l.Logger.Errorf("账号状态异常: %s, status=%d", in.Account, model.UserStatusDisabled)
		return nil, errors.New("账号已被禁用，请联系管理员")
	}

	// 4. 密码哈希验证
	ok, err := authx.VerifyPassword(in.Password, auth.Salt, auth.Password)
	if err != nil {
		l.Logger.Errorf("密码验证异常: %v, account=%s", err, in.Account)
		return nil, errors.New("登录失败，请重试")
	}
	if !ok {
		l.Logger.Errorf("密码验证失败: account=%s", in.Account)
		return nil, errors.New("账号或密码错误")
	}

	// 5. 生成JWT令牌
	accessToken, err := authx.GenJwtToken(
		l.svcCtx.Config.JWT.AccessSecret,
		l.svcCtx.Config.JWT.AccessExpire,
		map[string]any{
			"user_id": auth.UserID,
			"account": auth.Account,
		})

	if err != nil {
		l.Logger.Errorf("生成访问令牌失败: %v, account=%s", err, auth.Account)
		return nil, errors.New("登录失败，请重试")
	}

	refreshToken, err := authx.GenJwtToken(
		l.svcCtx.Config.JWT.RefreshSecret,
		l.svcCtx.Config.JWT.RefreshExpire,
		map[string]any{
			"user_id": auth.UserID,
			"account": auth.Account,
		})
	if err != nil {
		l.Logger.Errorf("生成刷新令牌失败: %v, account=%s", err, auth.Account)
		return nil, errors.New("登录失败，请重试")
	}

	// 单点登录和同平台设备限制逻辑
	// 平台类型默认为"unknown"
	platform := "unknown"
	if in.DeviceType != nil {
		platform = *in.DeviceType
	}

	// 构建Redis键：用户ID:平台类型
	redisKey := fmt.Sprintf("user:login:%d:%s", auth.UserID, platform)

	// 检查该用户在该平台是否已有登录记录
	oldToken, err := l.svcCtx.Redis.Get(redisKey)
	if err != nil && err != redis.Nil {
		l.Logger.Errorf("查询用户登录记录失败: %v, user_id=%d, platform=%s", err, auth.UserID, platform)
		return nil, errors.New("登录失败，请重试")
	}

	// 如果存在旧令牌，将其加入黑名单
	if oldToken != "" {
		// 设置旧令牌为黑名单，有效期与原令牌一致
		blacklistKey := fmt.Sprintf("token:blacklist:%s", oldToken)
		err = l.svcCtx.Redis.Set(blacklistKey, "1")
		if err != nil {
			l.Logger.Errorf("添加旧令牌到黑名单失败: %v, user_id=%d, platform=%s", err, auth.UserID, platform)
			return nil, errors.New("登录失败，请重试")
		}
		l.Logger.Infof("用户在同平台已有登录，已使旧令牌失效: user_id=%d, platform=%s", auth.UserID, platform)
	}

	// 存储新的登录信息到Redis
	err = l.svcCtx.Redis.Set(redisKey, accessToken)
	if err != nil {
		l.Logger.Errorf("存储用户登录记录失败: %v, user_id=%d, platform=%s", err, auth.UserID, platform)
		return nil, errors.New("登录失败，请重试")
	}

	// 6. todo 记录登录日志
	l.Logger.Infof("用户登录成功: user_id=%d, account=%s, device_type=%v",
		auth.UserID, in.Account, in.DeviceType)

	return &pb.LoginResp{
		AccessToken:        accessToken,
		AccessTokenExpire:  l.svcCtx.Config.JWT.AccessExpire,
		RefreshToken:       refreshToken,
		RefreshTokenExpire: l.svcCtx.Config.JWT.RefreshExpire,
	}, nil
}
