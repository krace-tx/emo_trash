package postlogic

import (
	"context"
	"time"

	"github.com/krace-tx/emo_trash/app/model"
	"github.com/krace-tx/emo_trash/app/rpc/post/internal/svc"
	"github.com/krace-tx/emo_trash/app/rpc/post/pb"
	errx "github.com/krace-tx/emo_trash/pkg/err"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/zeromicro/go-zero/core/logx"
)

type StarPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStarPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StarPostLogic {
	return &StarPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StarPostLogic) StarPost(in *pb.StarPostReq) (*pb.CommonResp, error) {
	postOid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, errx.ErrSystemArgInvalid
	}

	starColl := l.svcCtx.Mongo.Collection(model.StarCollectionName)
	postColl := l.svcCtx.Mongo.Collection(model.PostCollectionName)

	// 1. 检查帖子是否存在
	var p model.Post
	if err := postColl.FindOne(l.ctx, bson.M{"_id": postOid, "deleted_at": bson.M{"$exists": false}}).Decode(&p); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errx.ErrPostNotFound
		}
		return nil, errx.ErrDBQueryFailed
	}

	// 2. 切换收藏状态 (Toggle)
	filter := bson.M{"post_id": in.Id, "user_id": in.UserId}
	var existing model.Star
	err = starColl.FindOne(l.ctx, filter).Decode(&existing)

	if err == mongo.ErrNoDocuments {
		// 执行收藏
		_, err := starColl.InsertOne(l.ctx, &model.Star{
			PostId:    in.Id,
			UserId:    in.UserId,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return nil, errx.ErrDBInsertFailed
		}
		// 增加计数
		_, _ = postColl.UpdateOne(l.ctx, bson.M{"_id": postOid}, bson.M{"$inc": bson.M{"star_count": 1}})
		return &pb.CommonResp{Success: true, Message: "收藏成功"}, nil
	} else if err == nil {
		// 取消收藏
		_, err := starColl.DeleteOne(l.ctx, filter)
		if err != nil {
			return nil, errx.ErrDBDeleteFailed
		}
		// 减少计数
		_, _ = postColl.UpdateOne(l.ctx, bson.M{"_id": postOid}, bson.M{"$inc": bson.M{"star_count": -1}})
		return &pb.CommonResp{Success: true, Message: "取消收藏成功"}, nil
	}

	return nil, errx.ErrDBQueryFailed
}
