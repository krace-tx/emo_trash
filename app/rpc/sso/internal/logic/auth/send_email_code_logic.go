package authlogic

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"
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
	if !email.IsValidEmail(in.Email) {
		return nil, errx.ErrSystemArgInvalid
	}

	code, err := genVerifyCode(6)
	if err != nil {
		l.Logger.Errorf("生成邮箱验证码失败: %v, email=%s, scene=%s", err, in.Email, in.Scene)
		return nil, errx.ErrSystemInternal
	}

	codeKey := redis.GenerateKey("sso", "email_code", in.Scene, in.Email)
	expire := 5 * time.Minute
	if err = l.svcCtx.Redis.Set(codeKey, code, expire); err != nil {
		l.Logger.Errorf("缓存邮箱验证码失败: %v, email=%s, scene=%s", err, in.Email, in.Scene)
		return nil, errx.ErrSystemInternal
	}

	if err = email.SendCode(l.svcCtx.Config.Email, in.Email, code, in.Scene); err != nil {
		_ = l.svcCtx.Redis.Del(codeKey)
		l.Logger.Errorf("发送邮箱验证码失败: %v, email=%s, scene=%s", err, in.Email, in.Scene)
		return nil, errx.ErrSystemInternal
	}

	l.Logger.Infof("发送邮箱验证码成功: email=%s, scene=%s", in.Email, in.Scene)
	return &pb.SendEmailCodeResp{
		Success:       true,
		Message:       "验证码发送成功",
		ExpireSeconds: int32(expire / time.Second),
	}, nil
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
