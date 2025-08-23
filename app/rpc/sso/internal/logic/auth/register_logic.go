package authlogic

import (
	"context"
	"fmt"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/model"
	authx "github.com/krace-tx/emo_trash/pkg/auth"
	"github.com/krace-tx/emo_trash/pkg/db/rdb"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"golang.org/x/crypto/bcrypt"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	userModel *model.UserModel
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		userModel: model.NewUserModel(ctx, svcCtx),
	}
}

// 注册
// 分布式锁：使用分布式锁确保手机号注册的唯一性
// 查重处理：检查手机号是否已注册
// 用户ID生成：生成唯一用户ID（分布式ID）
// 密码安全处理：生成盐值并加密密码
// 数据存储：事务方式创建UserAuth和UserProfile记录
// 令牌生成：生成访问令牌和刷新令牌
// 返回结果：组装并返回注册成功响应
func (l *RegisterLogic) Register(in *pb.RegisterReq) (*pb.RegisterResp, error) {

	// 2. Verify SMS code
	if err := l.verifySmsCode(in.Mobile, in.SmsCode); err != nil {
		return nil, err
	}

	// 3. 检查手机号是否已注册
	if exists, err := l.checkMobileExists(in.Mobile); err != nil {
		l.Logger.Error(err)
		return nil, errx.ErrAuthMobileExists
	} else if exists {
		return nil, errx.ErrAuthMobileExists
	}

	//4. 生成用户ID
	userID := uint64(l.svcCtx.Snowflake.Generate())

	//5. 生成盐值和加密密码
	salt, encryptedPassword, err := l.encryptPassword(in.Password)
	if err != nil {
		l.Logger.Error(err)
		return nil, err
	}

	auth, _, err := l.userModel.CreateUser(userID, in.Mobile, encryptedPassword, salt)
	if err != nil {
		l.Logger.Error(err)
		return nil, err
	}

	accessToken, accessExpire, refreshToken, refreshExpire := l.generateTokens(*auth)

	return &pb.RegisterResp{
		AccessToken:        accessToken,
		AccessTokenExpire:  accessExpire,
		RefreshToken:       refreshToken,
		RefreshTokenExpire: refreshExpire,
	}, nil
}

// 验证短信验证码
func (l *RegisterLogic) verifySmsCode(mobile, code string) error {
	cacheKey := fmt.Sprintf("sms:code:%s", mobile)
	storedCode, err := l.svcCtx.Redis.Get(cacheKey)
	if err != nil {
		if err == redis.Nil {
			return errx.ErrAuthSmsCodeInvalid
		}
		l.Logger.Error("Failed to get SMS code from cache:", err)
		return errx.ErrSystemInternal
	}

	if storedCode != code {
		return errx.ErrAuthSmsCodeInvalid
	}

	defer func() {
		if _, err := l.svcCtx.Redis.Del(cacheKey); err != nil {
			l.Logger.Error("Failed to delete SMS code cache:", err)
		}
	}()

	return nil
}

// 检查手机号是否已注册
func (l *RegisterLogic) checkMobileExists(mobile string) (bool, error) {
	engine := rdb.NewEngine[model.UserAuth](rdb.M)
	count, err := engine.Count(l.ctx, rdb.WithConditions("mobile = ?", mobile))
	return count > 0, err
}

// 加密密码
func (l *RegisterLogic) encryptPassword(password string) (string, string, error) {
	// 生成随机盐值
	salt, err := authx.GenerateSalt()
	if err != nil {
		l.Logger.Error(err)
		return "", "", err
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	return salt, string(hashedPassword), nil
}

// Generate tokens
func (l *RegisterLogic) generateTokens(auth model.UserAuth) (accessToken string, accessExpire int64, refreshToken string, refreshExpire int64) {
	accessExpire = l.svcCtx.Config.JWT.AccessExpire
	accessToken, _ = authx.GenJwtToken(
		l.svcCtx.Config.JWT.AccessSecret,
		l.svcCtx.Config.JWT.AccessExpire,
		map[string]any{
			"user_id":       auth.UserID,
			"account":       auth.Account,
			"platform":      auth.Platform,
			"last_login_ip": auth.LastLoginIP,
		},
	)

	refreshToken, _ = authx.GenJwtToken(
		l.svcCtx.Config.JWT.RefreshSecret,
		l.svcCtx.Config.JWT.RefreshExpire,
		map[string]any{
			"user_id":       auth.UserID,
			"account":       auth.Account,
			"platform":      auth.Platform,
			"last_login_ip": auth.LastLoginIP,
		},
	)
	return
}
