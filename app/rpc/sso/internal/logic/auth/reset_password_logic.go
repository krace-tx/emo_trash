package authlogic

import (
	"context"
	"errors"
	"time"

	"github.com/krace-tx/emo_trash/app/model"
	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"
	authx "github.com/krace-tx/emo_trash/pkg/auth"
	consts "github.com/krace-tx/emo_trash/pkg/constant"
	"github.com/krace-tx/emo_trash/pkg/datastore/redis"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 忘记密码（通过邮箱验证码重置）
func (l *ResetPasswordLogic) ResetPassword(in *pb.ResetPasswordReq) (*pb.CommonResp, error) {
	// 1. 验证邮箱验证码
	if !l.verifyCode(in.Email, in.EmailCode, consts.SceneResetPwd) {
		return nil, errx.ErrAuthSmsCodeInvalid
	}

	// 2. 查询用户
	userColl := l.svcCtx.Mongo.Collection(model.UserCollectionName)
	filter := bson.M{
		"email":      in.Email,
		"deleted_at": bson.M{"$exists": false},
	}
	var user model.User
	if err := userColl.FindOne(l.ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			l.Logger.Errorf("重置密码用户不存在: email=%s", in.Email)
			return nil, errx.ErrAuthPasswordIncorrect
		}
		l.Logger.Errorf("查询重置密码用户失败: %v, email=%s", err, in.Email)
		return nil, errx.ErrDBQueryFailed
	}

	// 3. 加密新密码
	salt, err := authx.GenerateSalt()
	if err != nil {
		l.Logger.Errorf("生成重置密码盐失败: %v, email=%s", err, in.Email)
		return nil, errx.ErrAuthPwdEncryptFail
	}
	hashedPassword, err := authx.HashPassword(in.NewPassword, salt)
	if err != nil {
		l.Logger.Errorf("重置密码哈希失败: %v, email=%s", err, in.Email)
		return nil, errx.ErrAuthPwdEncryptFail
	}

	// 4. 更新用户密码
	update := bson.M{
		"$set": bson.M{
			"password":   hashedPassword,
			"salt":       salt,
			"updated_at": time.Now(),
		},
	}
	if _, err := userColl.UpdateOne(l.ctx, filter, update); err != nil {
		l.Logger.Errorf("重置密码更新失败: %v, email=%s", err, in.Email)
		return nil, errx.ErrDBUpdateFailed
	}

	// 5. 成功后清理验证码
	l.delCode(in.Email, consts.SceneResetPwd)

	l.Logger.Infof("密码重置成功: email=%s", in.Email)
	return &pb.CommonResp{Success: true, Message: "密码重置成功"}, nil
}

func (l *ResetPasswordLogic) verifyCode(email, code, scene string) bool {
	codeKey := redis.GenerateKey("sso", "email_code", scene, email)
	var cachedCode string
	if err := l.svcCtx.Redis.Get(codeKey, &cachedCode); err != nil {
		l.Logger.Errorf("获取验证码失败: %v, email=%s, scene=%s", err, email, scene)
		return false
	}
	return cachedCode == code
}

func (l *ResetPasswordLogic) delCode(email, scene string) {
	codeKey := redis.GenerateKey("sso", "email_code", scene, email)
	_ = l.svcCtx.Redis.Del(codeKey)
}
