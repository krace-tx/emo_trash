package authlogic

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"
	consts "github.com/krace-tx/emo_trash/pkg/constant"
	"github.com/krace-tx/emo_trash/pkg/datastore/redis"
	"github.com/krace-tx/emo_trash/pkg/email"
	errx "github.com/krace-tx/emo_trash/pkg/err"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendEmailCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendEmailCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendEmailCodeLogic {
	return &SendEmailCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 发送邮箱验证码
func (l *SendEmailCodeLogic) SendEmailCode(in *pb.SendEmailCodeReq) (*pb.SendEmailCodeResp, error) {
	// 1. 基础校验
	if !email.IsValidEmail(in.Email) {
		return nil, errx.ErrSystemArgInvalid
	}

	// 2. 场景校验
	if !l.isValidScene(in.Scene) {
		l.Logger.Errorf("非法场景参数: email=%s, scene=%s", in.Email, in.Scene)
		return nil, errx.ErrSystemArgInvalid
	}

	// 3. 频率校验 (30s 冷却)
	if l.isInCooling(in.Scene, in.Email) {
		return nil, errx.ErrAuthRateLimit
	}

	// 4. 生成验证码
	code, err := genVerifyCode(6)
	if err != nil {
		l.Logger.Errorf("生成邮箱验证码失败: %v, email=%s, scene=%s", err, in.Email, in.Scene)
		return nil, errx.ErrSystemInternal
	}

	// 5. 缓存验证码与冷却锁
	codeKey := redis.GenerateKey("sso", "email_code", in.Scene, in.Email)
	lockKey := redis.GenerateKey("sso", "email_code", "lock", in.Scene, in.Email)

	if err = l.svcCtx.Redis.Set(codeKey, code, consts.CodeExpireTime); err != nil {
		l.Logger.Errorf("缓存邮箱验证码失败: %v, email=%s, scene=%s", err, in.Email, in.Scene)
		return nil, errx.ErrSystemInternal
	}
	// 设置 30s 冷却锁
	_ = l.svcCtx.Redis.Set(lockKey, "1", consts.CodeCoolingTime)

	// 6. 发送邮件
	if err = email.SendCode(l.svcCtx.Config.Email, in.Email, code, in.Scene); err != nil {
		_ = l.svcCtx.Redis.Del(codeKey)
		_ = l.svcCtx.Redis.Del(lockKey)
		l.Logger.Errorf("发送邮箱验证码失败: %v, email=%s, scene=%s", err, in.Email, in.Scene)
		return nil, errx.ErrSystemInternal
	}

	l.Logger.Infof("发送邮箱验证码成功: email=%s, scene=%s", in.Email, in.Scene)
	return &pb.SendEmailCodeResp{
		Success:       true,
		Message:       "验证码发送成功",
		ExpireSeconds: int32(consts.CodeExpireTime / time.Second),
	}, nil
}

func (l *SendEmailCodeLogic) isValidScene(scene string) bool {
	switch scene {
	case consts.SceneRegister, consts.SceneLogin, consts.SceneResetPwd:
		return true
	default:
		return false
	}
}

func (l *SendEmailCodeLogic) isInCooling(scene, email string) bool {
	lockKey := redis.GenerateKey("sso", "email_code", "lock", scene, email)
	var locked string
	_ = l.svcCtx.Redis.Get(lockKey, &locked)
	return locked != ""
}

func genVerifyCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid verify code length: %d", length)
	}

	max := big.NewInt(10)
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = byte('0' + n.Int64())
	}
	return string(b), nil
}
