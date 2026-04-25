package authlogic

import (
	"context"
	"errors"
	"time"

	"github.com/krace-tx/emo_trash/app/model"
	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"
	authx "github.com/krace-tx/emo_trash/pkg/auth"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改密码（已登录状态，需旧密码验证）
func (l *ChangePasswordLogic) ChangePassword(in *pb.ChangePasswordReq) (*pb.CommonResp, error) {
	userId := in.UserId
	if userId == "" {
		return nil, errx.ErrAuthUnauthorized
	}

	uid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		l.Logger.Errorf("无效的用户ID: %v, user_id=%s", err, userId)
		return nil, errx.ErrSystemArgInvalid
	}

	// 2. 查询用户
	userColl := l.svcCtx.Mongo.Collection(model.UserCollectionName)
	filter := bson.M{"_id": uid, "deleted_at": bson.M{"$exists": false}}

	var user model.User
	if err := userColl.FindOne(l.ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errx.ErrUserNotFound
		}
		l.Logger.Errorf("查询用户失败: %v, user_id=%s", err, userId)
		return nil, errx.ErrDBQueryFailed
	}

	// 3. 验证旧密码
	ok, err := authx.VerifyPassword(in.OldPassword, user.Salt, user.Password)
	if err != nil {
		l.Logger.Errorf("旧密码验证异常: %v, user_id=%s", err, userId)
		return nil, errx.ErrAuthPasswordVerifyError
	}
	if !ok {
		return nil, errx.ErrAuthPasswordIncorrect
	}

	// 4. 加密新密码
	salt, err := authx.GenerateSalt()
	if err != nil {
		return nil, errx.ErrAuthPwdEncryptFail
	}
	hashedPassword, err := authx.HashPassword(in.NewPassword, salt)
	if err != nil {
		return nil, errx.ErrAuthPwdEncryptFail
	}

	// 5. 更新密码
	update := bson.M{
		"$set": bson.M{
			"password":   hashedPassword,
			"salt":       salt,
			"updated_at": time.Now(),
		},
	}
	if _, err := userColl.UpdateOne(l.ctx, filter, update); err != nil {
		l.Logger.Errorf("更新用户密码失败: %v, user_id=%s", err, userId)
		return nil, errx.ErrDBUpdateFailed
	}

	l.Logger.Infof("用户修改密码成功: user_id=%s", userId)
	return &pb.CommonResp{Success: true, Message: "密码修改成功"}, nil
}
