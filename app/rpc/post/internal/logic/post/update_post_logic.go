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

type UpdatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePostLogic {
	return &UpdatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePostLogic) UpdatePost(in *pb.UpdatePostReq) (*pb.CommonResp, error) {
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
		l.Logger.Errorf("越权修改帖子: post_id=%s, requester=%s, owner=%s", in.Id, in.UserId, p.AuthorId)
		return nil, errx.ErrPostUpdateAllowed
	}

	update := bson.M{
		"$set": bson.M{
			"title":      in.Title,
			"content":    in.Content,
			"images":     in.Images,
			"updated_at": time.Now(),
		},
	}
	if _, err := coll.UpdateOne(l.ctx, filter, update); err != nil {
		l.Logger.Errorf("更新帖子失败: %v, post_id=%s", err, in.Id)
		return nil, errx.ErrDBUpdateFailed
	}

	return &pb.CommonResp{Success: true, Message: "更新成功"}, nil
}
