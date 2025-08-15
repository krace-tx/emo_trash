package authlogic

import (
	"context"
	"errors"
	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/model"
	"github.com/krace-tx/emo_trash/pkg/db/rdb"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 注册
// 参数验证：验证手机号格式、密码强度、短信验证码有效性
// 查重处理：检查手机号是否已注册
// 用户ID生成：生成唯一用户ID（分布式ID）
// 密码安全处理：生成盐值并加密密码
// 数据存储：事务方式创建UserAuth和UserProfile记录
// 令牌生成：生成访问令牌和刷新令牌
// 返回结果：组装并返回注册成功响应
func (l *RegisterLogic) Register(in *pb.RegisterReq) (*pb.RegisterResp, error) {
	// 1. 参数验证
	if err := l.validateParams(in); err != nil {
		return nil, err
	}

	// 2. 验证短信验证码
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

	// 4. 生成用户ID
	//userID := l.generateUserID()

	// 5. 生成盐值和加密密码
	//salt, encryptedPassword := l.encryptPassword(in.Password)

	// 6. 事务创建用户数据
	//if err := l.createUser(userID, in.Mobile, encryptedPassword, salt); err != nil {
	//	return nil, err
	//}

	// 7. 生成访问令牌和刷新令牌
	//accessToken, accessExpire, refreshToken, refreshExpire := l.generateTokens(userID)

	// 8. 返回注册结果
	return &pb.RegisterResp{
		//AccessToken:        accessToken,
		//AccessTokenExpire:  accessExpire,
		//RefreshToken:       refreshToken,
		//RefreshTokenExpire: refreshExpire,
	}, nil
}

// 参数验证
func (l *RegisterLogic) validateParams(in *pb.RegisterReq) error {
	// 手机号格式验证
	if !regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(in.Mobile) {
		return errors.New("手机号格式不正确")
	}

	// 密码强度验证(8-32位，包含字母、数字和特殊字符)
	if len(in.Password) < 8 || len(in.Password) > 32 {
		return errors.New("密码长度必须为8-32位")
	}
	if !regexp.MustCompile("^(?=.*[a-zA-Z])(?=.*\\d)(?=.*[!@#$%^&*()_+{}|:\"<>?~`'\\\\[\\];\\',./]).+$").MatchString(in.Password) {
		return errors.New("密码必须包含字母、数字和特殊字符")
	}

	// 短信验证码验证
	if len(in.SmsCode) != 6 {
		return errors.New("短信验证码格式不正确")
	}

	return nil
}

// 验证短信验证码
func (l *RegisterLogic) verifySmsCode(mobile, code string) error {
	// 实际项目中需要调用短信服务验证验证码
	// 此处为示例，假设调用短信服务验证
	// if !smsService.VerifyCode(mobile, code) {
	//     return errors.New("短信验证码错误或已过期")
	// }
	return nil
}

// 检查手机号是否已注册
func (l *RegisterLogic) checkMobileExists(mobile string) (bool, error) {
	engine := rdb.NewEngine[model.UserAuth](rdb.M)
	count, err := engine.Count(l.ctx, rdb.WithConditions("mobile = ?", mobile))
	return count > 0, err
}

// 生成用户ID(使用雪花算法或其他分布式ID生成策略)
func (l *RegisterLogic) generateUserID() uint64 {
	// TODO 分布式ID生成器
	// 此处为示例，使用当前时间戳+随机数简单生成
	return uint64(time.Now().UnixNano()/1e6) + rand.Uint64()
}

// 加密密码
//func (l *RegisterLogic) encryptPassword(password string) (string, string) {
//	// 生成随机盐值
//	salt := generateRandomSalt(16)
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
//	return salt, string(hashedPassword)
//}

// 创建用户数据(事务)
func (l *RegisterLogic) createUser(userID uint64, mobile, password, salt string) error {
	transfunc := func(engine rdb.EngineInterface[model.UserAuth]) error {
		auth := &model.UserAuth{
			UserID:   userID,
			Mobile:   mobile,
			Account:  mobile,
			Password: password,
			Salt:     salt,
			Status:   model.UserStatusNormal,
		}
		if err := engine.Create(l.ctx, auth); err != nil {
			return errx.ErrAuthCreateUserAuthFail
		}

		profileEngine := rdb.NewEngine[model.UserProfile](rdb.M)
		profile := &model.UserProfile{
			UserID:   userID,
			Nickname: "用户" + strconv.FormatUint(userID, 10)[8:],
			Avatar:   "",
		}
		if err := profileEngine.Create(l.ctx, profile); err != nil {
			return errx.ErrAuthCreateUserProfileFail
		}

		return nil
	}

	return rdb.Transaction[model.UserAuth](l.ctx, rdb.M, transfunc)
}
