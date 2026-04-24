package authlogic

import (
	"context"
	"errors"
	"time"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/model"
	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"
	authx "github.com/krace-tx/emo_trash/pkg/auth"
	consts "github.com/krace-tx/emo_trash/pkg/constant"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

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

// 邮箱注册
func (l *RegisterLogic) Register(in *pb.RegisterReq) (*pb.LoginResp, error) {
	userColl := l.svcCtx.Mongo.Collection("users")
	existFilter := bson.M{
		"email":      in.Email,
		"deleted_at": bson.M{"$exists": false},
	}

	// 校验邮箱是否已注册
	existsErr := userColl.FindOne(l.ctx, existFilter).Err()
	if existsErr == nil {
		return nil, errx.ErrAuthMobileExists
	}
	if !errors.Is(existsErr, mongo.ErrNoDocuments) {
		l.Logger.Errorf("查询注册邮箱失败: %v, email=%s", existsErr, in.Email)
		return nil, errx.ErrDBQueryFailed
	}

	// 密码加盐哈希
	salt, err := authx.GenerateSalt()
	if err != nil {
		l.Logger.Errorf("生成密码盐失败: %v, email=%s", err, in.Email)
		return nil, errx.ErrAuthPwdEncryptFail
	}
	hashedPassword, err := authx.HashPassword(in.Password, salt)
	if err != nil {
		l.Logger.Errorf("密码哈希失败: %v, email=%s", err, in.Email)
		return nil, errx.ErrAuthPwdEncryptFail
	}

	now := time.Now()
	user := &model.User{
		Email:     in.Email,
		Salt:      salt,
		Password:  hashedPassword,
		Nickname:  in.Email,
		Avatar:    "",
		Status:    model.UserStatusNormal,
		CreatedAt: now,
		UpdatedAt: now,
	}

	insertResult, err := userColl.InsertOne(l.ctx, user)
	if err != nil {
		l.Logger.Errorf("创建用户失败: %v, email=%s", err, in.Email)
		return nil, errx.ErrAuthRegisterFailed
	}

	insertedID, ok := insertResult.InsertedID.(interface{ Hex() string })
	if !ok {
		l.Logger.Errorf("获取新用户ID失败: inserted_id=%v", insertResult.InsertedID)
		return nil, errx.ErrAuthGenIDFailed
	}

	claims := map[string]any{
		consts.UserId: insertedID.Hex(),
		"email":       in.Email,
	}

	accessToken, err := authx.GenJwtToken(
		l.svcCtx.Config.JWT.AccessSecret,
		l.svcCtx.Config.JWT.AccessExpire,
		claims,
	)
	if err != nil {
		l.Logger.Errorf("注册后生成访问令牌失败: %v, email=%s", err, in.Email)
		return nil, errx.ErrAuthGenAccessTokenFail
	}

	refreshToken, err := authx.GenJwtToken(
		l.svcCtx.Config.JWT.RefreshSecret,
		l.svcCtx.Config.JWT.RefreshExpire,
		claims,
	)
	if err != nil {
		l.Logger.Errorf("注册后生成刷新令牌失败: %v, email=%s", err, in.Email)
		return nil, errx.ErrAuthGenRefreshTokenFail
	}

	l.Logger.Infof("用户注册成功: email=%s, user_id=%s", in.Email, insertedID.Hex())
	return &pb.LoginResp{
		AccessToken:        accessToken,
		AccessTokenExpire:  l.svcCtx.Config.JWT.AccessExpire,
		RefreshToken:       refreshToken,
		RefreshTokenExpire: l.svcCtx.Config.JWT.RefreshExpire,
	}, nil
}
