package authlogic

import (
	"context"
	"errors"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/model"
	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户信息
func (l *GetUserInfoLogic) GetUserInfo(in *pb.GetUserInfoReq) (*pb.GetUserInfoResp, error) {
	uid, err := primitive.ObjectIDFromHex(in.UserId)
	if err != nil {
		l.Logger.Errorf("无效的用户ID: %v, user_id=%s", err, in.UserId)
		return nil, errx.ErrSystemArgInvalid
	}

	userColl := l.svcCtx.Mongo.Collection(model.UserCollectionName)
	filter := bson.M{
		"_id":        uid,
		"deleted_at": bson.M{"$exists": false},
	}

	var user model.User
	if err := userColl.FindOne(l.ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errx.ErrUserNotFound
		}
		l.Logger.Errorf("查询用户信息失败: %v, user_id=%s", err, in.UserId)
		return nil, errx.ErrDBQueryFailed
	}

	return &pb.GetUserInfoResp{
		User: &pb.UserInfo{
			UserId:     user.ID.Hex(),
			Email:      user.Email,
			Nickname:   user.Nickname,
			Avatar:     user.Avatar,
			CreateTime: user.CreatedAt.Unix(),
			Bio:        user.Bio,
			Mood:       user.Mood,
		},
	}, nil
}
