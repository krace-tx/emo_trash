package authlogic

import (
	"context"
	"errors"

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

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 邮箱登录
func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {
	userColl := l.svcCtx.Mongo.Collection(model.UserCollectionName)
	filter := bson.M{
		"email":      in.Email,
		"deleted_at": bson.M{"$exists": false},
	}

	var user model.User
	if err := userColl.FindOne(l.ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errx.ErrAuthPasswordIncorrect
		}
		l.Logger.Errorf("查询登录用户失败: %v, email=%s", err, in.Email)
		return nil, errx.ErrDBQueryFailed
	}

	if user.Status == model.UserStatusDisabled {
		l.Logger.Errorf("登录账号已禁用: email=%s", in.Email)
		return nil, errx.ErrUserDisabled
	}

	ok, err := authx.VerifyPassword(in.Password, user.Salt, user.Password)
	if err != nil {
		l.Logger.Errorf("密码验证异常: %v, email=%s", err, in.Email)
		return nil, errx.ErrAuthPasswordVerifyError
	}
	if !ok {
		l.Logger.Errorf("密码验证失败: email=%s", in.Email)
		return nil, errx.ErrAuthPasswordIncorrect
	}

	l.Logger.Infof("用户登录成功: email=%s, user_id=%s", user.Email, user.ID.Hex())
	return l.generateTokenPair(user.ID.Hex(), user.Email)
}

func (l *LoginLogic) generateTokenPair(userId, email string) (*pb.LoginResp, error) {
	claims := map[string]any{
		consts.UserId: userId,
		"email":       email,
	}

	accessToken, err := authx.GenJwtToken(
		l.svcCtx.Config.JWT.AccessSecret,
		l.svcCtx.Config.JWT.AccessExpire,
		claims,
	)
	if err != nil {
		l.Logger.Errorf("生成访问令牌失败: %v, user_id=%s", err, userId)
		return nil, errx.ErrAuthGenAccessTokenFail
	}

	refreshToken, err := authx.GenJwtToken(
		l.svcCtx.Config.JWT.RefreshSecret,
		l.svcCtx.Config.JWT.RefreshExpire,
		claims,
	)
	if err != nil {
		l.Logger.Errorf("生成刷新令牌失败: %v, user_id=%s", err, userId)
		return nil, errx.ErrAuthGenRefreshTokenFail
	}

	return &pb.LoginResp{
		AccessToken:        accessToken,
		AccessTokenExpire:  l.svcCtx.Config.JWT.AccessExpire,
		RefreshToken:       refreshToken,
		RefreshTokenExpire: l.svcCtx.Config.JWT.RefreshExpire,
	}, nil
}
