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

type LikePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikePostLogic {
	return &LikePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LikePostLogic) LikePost(in *pb.LikePostReq) (*pb.CommonResp, error) {
	postOid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, errx.ErrSystemArgInvalid
	}

	likeColl := l.svcCtx.Mongo.Collection(model.LikeCollectionName)
	postColl := l.svcCtx.Mongo.Collection(model.PostCollectionName)

	// 1. 检查帖子是否存在
	var p model.Post
	if err := postColl.FindOne(l.ctx, bson.M{"_id": postOid, "deleted_at": bson.M{"$exists": false}}).Decode(&p); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errx.ErrPostNotFound
		}
		return nil, errx.ErrDBQueryFailed
	}

	// 2. 切换点赞状态 (Toggle)
	filter := bson.M{"post_id": in.Id, "user_id": in.UserId}
	var existing model.Like
	err = likeColl.FindOne(l.ctx, filter).Decode(&existing)

	if err == mongo.ErrNoDocuments {
		// 执行点赞
		_, err := likeColl.InsertOne(l.ctx, &model.Like{
			PostId:    in.Id,
			UserId:    in.UserId,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return nil, errx.ErrDBInsertFailed
		}
		// 增加计数
		_, _ = postColl.UpdateOne(l.ctx, bson.M{"_id": postOid}, bson.M{"$inc": bson.M{"like_count": 1}})
		return &pb.CommonResp{Success: true, Message: "点赞成功"}, nil
	} else if err == nil {
		// 取消点赞
		_, err := likeColl.DeleteOne(l.ctx, filter)
		if err != nil {
			return nil, errx.ErrDBDeleteFailed
		}
		// 减少计数
		_, _ = postColl.UpdateOne(l.ctx, bson.M{"_id": postOid}, bson.M{"$inc": bson.M{"like_count": -1}})
		return &pb.CommonResp{Success: true, Message: "取消点赞成功"}, nil
	}

	return nil, errx.ErrDBQueryFailed
}
