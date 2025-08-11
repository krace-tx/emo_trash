package authlogic

import (
	"context"
	"fmt" // 新增导入

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/model"
	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"
	authx "github.com/krace-tx/emo_trash/pkg/auth"
	"github.com/krace-tx/emo_trash/pkg/db/rdb"
	errx "github.com/krace-tx/emo_trash/pkg/err"
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

func (l *LoginByPasswordLogic) LoginByPassword(in *pb.LoginByPasswordReq) (*pb.LoginResp, error) {
	// 1. 基础参数校验
	if in.Account == "" || in.Password == "" {
		return nil, errx.ErrSystemArgInvalid
	}

	// 2. 查询用户认证信息
	engine := rdb.NewEngine[model.UserAuth](rdb.M)
	auth, err := engine.GetByCondition(l.ctx, rdb.WithConditions("account = ?", in.Account))
	if err != nil {
		l.Logger.Errorf("查询用户认证信息失败: %v, account=%s", err, in.Account)
		return nil, errx.ErrDBQueryFailed
	}
	if auth == nil {
		return nil, errx.ErrUserNotFound
	}

	// 3. 检查用户状态
	if auth.Status == model.UserStatusDisabled {
		l.Logger.Errorf("账号状态异常: %s, status=%d", in.Account, model.UserStatusDisabled)
		return nil, errx.ErrUserDisabled
	}

	// 4. 密码哈希验证
	ok, err := authx.VerifyPassword(in.Password, auth.Salt, auth.Password)
	if err != nil {
		l.Logger.Errorf("密码验证异常: %v, account=%s", err, in.Account)
		return nil, errx.ErrAuthPasswordVerifyError
	}
	if !ok {
		l.Logger.Errorf("密码验证失败: account=%s", in.Account)
		return nil, errx.ErrAuthPasswordIncorrect
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
		return nil, errx.ErrAuthGenAccessTokenFail
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
		return nil, errx.ErrAuthGenRefreshTokenFail
	}

	// 单点登录和同平台设备限制逻辑
	platform := "unknown"
	if in.DeviceType != nil {
		platform = *in.DeviceType
	}

	redisKey := fmt.Sprintf("user:login:%d:%s", auth.UserID, platform)

	oldToken, err := l.svcCtx.Redis.Get(redisKey)
	if err != nil && err != redis.Nil {
		l.Logger.Errorf("查询用户登录记录失败: %v, user_id=%d, platform=%s", err, auth.UserID, platform)
		return nil, errx.ErrAuthSSOCheckFail
	}

	if oldToken != "" {
		blacklistKey := fmt.Sprintf("token:blacklist:%s", oldToken)
		err = l.svcCtx.Redis.Set(blacklistKey, "1")
		if err != nil {
			l.Logger.Errorf("添加旧令牌到黑名单失败: %v, user_id=%d, platform=%s", err, auth.UserID, platform)
			return nil, errx.ErrAuthTokenBlacklistFail
		}
		l.Logger.Infof("用户在同平台已有登录，已使旧令牌失效: user_id=%d, platform=%s", auth.UserID, platform)
	}

	if err := l.svcCtx.Redis.Set(redisKey, accessToken); err != nil {
		l.Logger.Errorf("存储用户登录记录失败: %v, user_id=%d, platform=%s", err, auth.UserID, platform)
		return nil, errx.ErrAuthSaveLoginRecordFail
	}

	// 6. 记录登录日志
	//loginLog := &model.LoginLog{
	//	UserID:     auth.UserID, // 假设auth.UserID为primitive.ObjectID类型
	//	Account:    in.Account,
	//	Platform:   platform,            // 已在前面定义的平台标识(web/ios/android等)
	//	DeviceID:   *in.DeviceId,        // 从请求中获取设备ID
	//	IP:         *in.LoginIp,  // 获取客户端IP
	//	UserAgent:  getUserAgent(l.ctx), // 获取用户代理信息
	//	LoginTime:  time.Now(),
	//	Status:     "success", // 登录状态：成功
	//	ResultCode: 200,       // 成功状态码
	//	ResultMsg:  "登录成功",
	//	TraceID:    logx.TraceIDFromContext(l.ctx), // 从上下文获取追踪ID
	//	Extra: map[string]interface{}{
	//		"login_type": "password", // 登录方式：密码登录
	//		"client":     in.Client,  // 客户端标识(假设请求中有该字段)
	//	},
	//}
	//
	//// 插入登录日志到MongoDB
	//_, err = l.svcCtx.Mongo.Collection("login_logs").InsertOne(l.ctx, loginLog)
	//if err != nil {
	//	// 日志记录失败不影响登录主流程，但需记录错误
	//	l.Logger.Errorf("记录登录日志失败: %v, user_id=%s, account=%s", err, auth.UserID, in.Account)
	//}
	//
	//l.Logger.Infof("用户登录成功: user_id=%s, account=%s, device_type=%v",
	//	auth.UserID, in.Account, in.DeviceType)

	return &pb.LoginResp{
		AccessToken:        accessToken,
		AccessTokenExpire:  l.svcCtx.Config.JWT.AccessExpire,
		RefreshToken:       refreshToken,
		RefreshTokenExpire: l.svcCtx.Config.JWT.RefreshExpire,
	}, nil
}
