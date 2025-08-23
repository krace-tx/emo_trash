package authlogic

import (
	"context"
	"time"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/model"
	authx "github.com/krace-tx/emo_trash/pkg/auth"
	"github.com/krace-tx/emo_trash/pkg/db/no_sql"
	"github.com/krace-tx/emo_trash/pkg/db/rdb"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	userModel  *model.UserModel
	tokenModel *model.TokenModel
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		userModel:  model.NewUserModel(ctx, svcCtx),
		tokenModel: model.NewTokenModel(ctx, svcCtx),
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
	lockKey := no_sql.GenerateKey("register", in.Mobile)
	lockID, ok, err := l.svcCtx.Redis.LockWithOptions(lockKey, 10*time.Second)
	if err != nil {
		return nil, errx.ErrSystemInternal
	}
	if !ok {
		return nil, errx.ErrAuthMobileExists
	}
	defer func() {
		if err := l.svcCtx.Redis.UnlockWithID(lockKey, lockID); err != nil {
			logx.Error(err)
			return
		}
	}()

	// 2. Verify SMS code
	if err := l.verifySmsCode(in.Mobile, in.SmsCode); err != nil {
		//return nil, err
	}

	// 3. 检查手机号是否已注册
	if exists, err := l.checkMobileExists(in.Mobile); err != nil {
		l.Logger.Error(err)
		return nil, err
	} else if exists {
		return nil, errx.ErrAuthMobileExists
	}

	//5. 生成盐值和加密密码
	salt, encryptedPassword, err := l.encryptPassword(in.Password)
	if err != nil {
		l.Logger.Error(err)
		return nil, err
	}

	auth, _, err := l.userModel.CreateUser(in.Account, in.Mobile, encryptedPassword, salt)
	if err != nil {
		l.Logger.Error(err)
		return nil, err
	}

	token, err := l.tokenModel.GenerateTokens(*auth)
	if err != nil {
		l.Logger.Error(err)
		return nil, err
	}

	return &pb.RegisterResp{
		AccessToken:        token.AccessToken,
		AccessTokenExpire:  token.AccessExpire,
		RefreshToken:       token.RefreshToken,
		RefreshTokenExpire: token.RefreshExpire,
	}, nil
}

// 验证短信验证码
func (l *RegisterLogic) verifySmsCode(mobile, code string) error {
	// 从缓存中获取短信验证码
	cacheKey := no_sql.GenerateKey("sms:code", mobile)
	var cacheCode string
	err := l.svcCtx.Redis.Get(cacheKey, &cacheCode)
	if err != nil {
		if err == redis.Nil {
			return errx.ErrAuthSmsCodeInvalid
		}
		l.Logger.Error("Failed to get SMS code from cache:", err)
		return errx.ErrSystemInternal
	}

	// 比较验证码
	if cacheCode != code {
		return errx.ErrAuthSmsCodeInvalid
	}

	// 删除缓存中的验证码
	if err := l.svcCtx.Redis.Del(cacheKey); err != nil {
		l.Logger.Error("Failed to delete SMS code cache:", err)
	}

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
