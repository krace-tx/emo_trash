package authlogic

import (
	"context"
	"time"

	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/model"
	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新用户信息
func (l *UpdateUserInfoLogic) UpdateUserInfo(in *pb.UpdateUserInfoReq) (*pb.CommonResp, error) {
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

	updateFields := bson.M{
		"updated_at": time.Now(),
	}

	if in.Nickname != "" {
		updateFields["nickname"] = in.Nickname
	}
	if in.Avatar != "" {
		updateFields["avatar"] = in.Avatar
	}
	if in.Bio != "" {
		updateFields["bio"] = in.Bio
	}
	if in.Mood != "" {
		updateFields["mood"] = in.Mood
	}

	update := bson.M{
		"$set": updateFields,
	}

	result, err := userColl.UpdateOne(l.ctx, filter, update)
	if err != nil {
		l.Logger.Errorf("更新用户信息失败: %v, user_id=%s", err, in.UserId)
		return nil, errx.ErrDBUpdateFailed
	}

	if result.MatchedCount == 0 {
		return nil, errx.ErrUserNotFound
	}

	return &pb.CommonResp{
		Success: true,
		Message: "用户信息更新成功",
	}, nil
}
