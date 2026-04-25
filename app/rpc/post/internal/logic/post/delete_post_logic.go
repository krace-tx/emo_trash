package postlogic

import (
	"context"
	"errors"
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

type DeletePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePostLogic {
	return &DeletePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeletePostLogic) DeletePost(in *pb.DeletePostReq) (*pb.CommonResp, error) {
	oid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, errx.ErrSystemArgInvalid
	}

	coll := l.svcCtx.Mongo.Collection(model.PostCollectionName)
	filter := bson.M{"_id": oid, "deleted_at": bson.M{"$exists": false}}

	var p model.Post
	if err := coll.FindOne(l.ctx, filter).Decode(&p); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errx.ErrPostNotFound
		}
		return nil, errx.ErrDBQueryFailed
	}

	if p.AuthorId != in.UserId {
		l.Logger.Errorf("越权删除帖子: post_id=%s, requester=%s", in.Id, in.UserId)
		return nil, errx.ErrPostDeleteAllowed
	}

	// 软删除
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
			"status":     -1,
		},
	}
	if _, err := coll.UpdateOne(l.ctx, filter, update); err != nil {
		return nil, errx.ErrDBUpdateFailed
	}

	return &pb.CommonResp{Success: true, Message: "删除成功"}, nil
}
