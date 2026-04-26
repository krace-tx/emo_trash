package authlogic

import (
	"context"
	"time"

	"github.com/krace-tx/emo_trash/app/model"
	"github.com/krace-tx/emo_trash/app/rpc/sso/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/sso/pb"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserStatsLogic {
	return &GetUserStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserStats 获取用户统计数据 (对应 mine.uvue 主卡片)
func (l *GetUserStatsLogic) GetUserStats(in *pb.GetUserStatsReq) (*pb.GetUserStatsResp, error) {
	l.Logger.Infof("GetUserStats user_id: %s", in.UserId)

	uid, err := primitive.ObjectIDFromHex(in.UserId)
	if err != nil {
		l.Logger.Errorf("非法用户ID: %s", in.UserId)
		return nil, errx.ErrSystemArgInvalid
	}

	// 1. 获取用户信息计算入驻天数
	userColl := l.svcCtx.Mongo.Collection(model.UserCollectionName)
	var user model.User
	if err := userColl.FindOne(l.ctx, bson.M{"_id": uid}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errx.ErrUserNotFound
		}
		l.Logger.Errorf("查询用户统计失败(User): %v", err)
		return nil, errx.ErrDBQueryFailed
	}

	joinDays := int64(time.Since(user.CreatedAt).Hours()/24) + 1

	// 2. 从 Mongo 获取帖子数和共鸣数
	postColl := l.svcCtx.Mongo.Collection(model.PostCollectionName)

	postCount, err := postColl.CountDocuments(l.ctx, bson.M{
		"author_id":  in.UserId,
		"deleted_at": bson.M{"$exists": false},
	})
	if err != nil {
		l.Logger.Errorf("查询用户统计失败(PostCount): %v", err)
		return nil, errx.ErrDBQueryFailed
	}

	// 计算点赞总数
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"author_id":  in.UserId,
			"deleted_at": bson.M{"$exists": false},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":        nil,
			"totalLikes": bson.M{"$sum": "$like_count"},
		}}},
	}
	cursor, err := postColl.Aggregate(l.ctx, pipeline)
	var resonanceCount int64 = 0
	if err == nil {
		var results []bson.M
		if err = cursor.All(l.ctx, &results); err == nil && len(results) > 0 {
			if val, ok := results[0]["totalLikes"].(int64); ok {
				resonanceCount = val
			} else if val, ok := results[0]["totalLikes"].(int32); ok {
				resonanceCount = int64(val)
			}
		}
	}

	return &pb.GetUserStatsResp{
		PostCount:      postCount,
		ResonanceCount: resonanceCount,
		JoinDays:       joinDays,
	}, nil
}
